package controllers

import (
	"fmt"
)

func getDerivedReadings() string {
	stmt := fmt.Sprintln(`
		derived_readings as (
			select
				instrument_id,
				sensor_id,
				value,
				read_at,
				created_at,
				deleted_at
			from
				readings r
			where
			(true)
			order by
				r.read_at asc
		)
	`)
	return stmt
}

func getCteReadings() string {
	stmt := fmt.Sprintln(`cte_readings as (
		select
			s.instrument_id,
			r.sensor_id,
			s."label",
			json_agg(
				json_build_object(
					'value', concat(r.value, s.unit_of_measurement),
					'read_at', r.read_at,
					'created_at',r.created_at
				)::json
			)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		where
			r.deleted_at is null
		group by
			s.instrument_id,
			r.sensor_id,
			s."label"
	)`)
	return stmt
}

//Required single params string;
//ex: 2021-01-23
func getCteReadingByHours() string {
	stmt := fmt.Sprintln(`
	hourly_series as (select generate_series(?::timestamp, ?::timestamp, '1 hour'::interval) as hour)
        , derived_readings as (
                select 
                        s.sensor_id, 
                        s.label, 
                        s.unit_of_measurement, 
                        s.measure, 
                        hs.hour as read_at, 
                        coalesce(round(avg(r.value::numeric), 3),0) as value
                from sensors s
                cross join hourly_series hs
                left join readings r 
                        on r.sensor_id = s.sensor_id and hs.hour = date_trunc('hour', r.read_at) and date(r.read_at) = ?
                        group by 1,2,3,4,5
        )
        ,cte_readings as (
                select
                                                s.instrument_id,
                                                s.measure,
                                                s.label,
                                                s.sensor_id,
                                                s.unit_of_measurement,
                                                json_agg(
                                                                json_build_object(
                                                                                                'read_at', r.read_at,
                                                                                                'value', r.value
                                                                )::json
                                                        order by r.read_at)::json as values
                from derived_readings r
                left join sensors s on s.sensor_id = r.sensor_id
                group by
                                                s.instrument_id,
                                                s.measure,
                                                s.label,
                                                s.sensor_id,
                                                s.unit_of_measurement
)`)
	return stmt
}

//Required single params string;
//ex: 2021-01-23
func getCteReadingByHoursV2() string {
	stmt := fmt.Sprintln(`
	derived_readings as (
		select
			r.sensor_id,
			extract('hour' from  date_trunc('hour', r.read_at)) as read_at,
			round(avg(r.value::numeric), 3) as value
		from
			readings r
		where
			date(r.read_at) = ?
		group by
			1, 2
	)
	,cte_readings as (
		select
				s.instrument_id,
				r.sensor_id,
				s."label",
				s.measure,
				json_agg(
						json_build_object(
								'value', r.value,
								'uom', s.unit_of_measurement,
								'read_at', r.read_at
						)::json
						order by r.read_at)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
				s.instrument_id,
				r.sensor_id,
				s."label",
				s.measure
)`)
	return stmt
}

func getCteReadingByHoursV3() string {
	stmt := fmt.Sprintln(`
	derived_readings as (
		select
				r.sensor_id,
				extract('hour' from  date_trunc('hour', r.read_at)) as read_at,
				round(avg(r.value::numeric), 3) as value
		from
				readings r
		where
				date(r.read_at) = ?
		group by
				1, 2
		order by
				2 asc, 3
)
,cte_readings as (
		select
						s.instrument_id,
						s.measure,
						r.read_at,
						json_agg(
										json_build_object(
														'sensor_id', r.sensor_id,
														'label', s.label,
														'value', r.value,
														'uom', s.unit_of_measurement
										)::json
						)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
						s.instrument_id,
						s.measure,
						r.read_at
)`)
	return stmt
}

func getCteReadingByDays() string {
	stmt := fmt.Sprintln(`
	daily_series as (select generate_series(?::timestamp, ?::timestamp, '1 day'::interval) as day)
	, derived_readings as (
		select 
			s.sensor_id, 
			s.label, 
			s.unit_of_measurement, 
			s.measure, 
			ds.day as read_at, 
			coalesce(round(avg(r.value::numeric), 3),0) as value
		from sensors s
		cross join daily_series ds
		left join readings r 
			on r.sensor_id = s.sensor_id and ds.day = date_trunc('day', r.read_at) and date(r.read_at) between ? and ?
			group by 1,2,3,4,5
	)
	,cte_readings as (
		select
										s.instrument_id,
										s.measure,
										s.label,
										s.sensor_id,
										s.unit_of_measurement,
										json_agg(
														json_build_object(
																						'read_at', r.read_at,
																						'value', r.value
														)::json
												order by r.read_at)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
										s.instrument_id,
										s.measure,
										s.label,
										s.sensor_id,
										s.unit_of_measurement
)
	`)
	return stmt
}

func getCteReadingByDaysV2() string {
	stmt := fmt.Sprintln(`
	derived_readings as (
		select 
			r.sensor_id,
			date(r.read_at) as read_at,
			round(avg(r.value::numeric),3) as value
		from readings r
		where
			date(r.read_at) between ? and ?
		group by 
			1, 2
		order by 2 asc
	)
	,cte_readings as (
		select
						s.instrument_id,
						s.measure,
						r.read_at,
						json_agg(
										json_build_object(
														'sensor_id', r.sensor_id,
														'label', s.label,
														'value', r.value,
														'uom', s.unit_of_measurement
										)::json
						)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
						s.instrument_id,
						s.measure,
						r.read_at
)
	`)
	return stmt
}

/*
	required single params of year
*/
func getCteReadingByMonthsOfYear() string {
	stmt := fmt.Sprintln(`
	monthly_series as (
		select generate_series(?::date, ?::date, '1 month'::interval) as months
	  ), derived_readings as (
		select 
		  s.sensor_id, 
		  s.label, 
		  s.unit_of_measurement, 
		  s.measure, 
		  ms.months as read_at, 
		  coalesce(round(avg(r.value::numeric), 3), 0) as value
		from 
		  sensors s
		cross join 
		  monthly_series ms
		left join 
		  readings r on r.sensor_id = s.sensor_id 
					 and ms.months = date_trunc('month', r.read_at) 
					 and extract(year from r.read_at) = ?
		group by 
		  1, 2, 3, 4, 5
	  )
        ,cte_readings as (
                select
					s.instrument_id,
					s.measure,
					s.label,
					s.sensor_id,
					s.unit_of_measurement,
					json_agg(
							json_build_object(
								'read_at', to_char(r.read_at, 'mon'),
								'value', r.value
							)::json
						order by r.read_at asc)::json as values
                from derived_readings r
                left join sensors s on s.sensor_id = r.sensor_id
                group by
                                                s.instrument_id,
                                                s.measure,
                                                s.label,
                                                s.sensor_id,
                                                s.unit_of_measurement
)`)
	return stmt
}

func getCteReadingByMonthsOfYearV2() string {
	stmt := fmt.Sprintln(`
	derived_readings as (
		select
			r.sensor_id,
			to_char(date(r.read_at), 'mon') as read_at,
			round(avg(r.value::numeric), 3) as value
		from
			readings r
		where
			date_part('year', r.read_at) = ?
		group by
			1, 2, extract(month from date(r.read_at))
		order by
			extract(month from date(r.read_at))
	)
	,cte_readings as (
		select
						s.instrument_id,
						s.measure,
						r.read_at,
						json_agg(
										json_build_object(
														'sensor_id', r.sensor_id,
														'label', s.label,
														'value', r.value,
														'uom', s.unit_of_measurement
										)::json
						)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
						s.instrument_id,
						s.measure,
						r.read_at
)
	`)
	return stmt
}

func getCteInstrument() string {
	stmt := fmt.Sprintln(`cte_instruments as (
		select
			i.incinerator_id,
			i.instrument_id,
			i.instrument_name,
			i.instrument_code,
			json_agg(
				json_build_object(
					'sensor_id', cte_r.sensor_id,
					'sensor_label', cte_r.label, 
					'sensor_values', cte_r.values
				)::json
			)::json as sensors
		from instruments i
		left join cte_readings cte_r on i.instrument_id = cte_r.instrument_id
		where i.incinerator_id = ? and i.instrument_id = ? and i.deleted_at is null
		group by
			i.incinerator_id,
			i.instrument_id,
			i.instrument_name,
			i.instrument_code
	)`)
	return stmt
}

func getCteInstrumentV2() string {
	stmt := fmt.Sprintln(`cte_instruments as (
		with grouped_measurement as (
				select
							i.incinerator_id,
							i.instrument_id,
							i.instrument_name,
							i.instrument_code,
							cte_r.measure,
							json_agg(
									json_build_object(
											'sensor_id', cte_r.sensor_id,
											'sensor_label', cte_r.label, 
											'sensor_values', cte_r.values
									)::json
							)::json as sensors
					from instruments i
					left join cte_readings cte_r on i.instrument_id = cte_r.instrument_id
					where
							i.deleted_at is null
					group by
							i.incinerator_id,
							i.instrument_id,
							i.instrument_name,
							i.instrument_code,
							cte_r.measure
			)
		select 
			gm.incinerator_id,
			gm.instrument_id,
			gm.instrument_name,
			gm.instrument_code,
			json_agg(
				json_build_object(
					'measure', gm.measure,
					'sensors', gm.sensors
				)::json
			)::json as measurement
		from grouped_measurement as gm
		group by
			gm.incinerator_id,
			gm.instrument_id,
			gm.instrument_name,
			gm.instrument_code
        )`)
	return stmt
}

func getCteInstrumentV3() string {
	stmt := fmt.Sprintln(`cte_instruments as (
		with grouped_measurement as (
						select
												i.incinerator_id,
												i.instrument_id,
												i.instrument_name,
												i.instrument_code,
												cte_r.measure,
												json_agg(
																json_build_object(
																				'read_at', cte_r.read_at,
																				'sensor_values', cte_r.values
																)::json
												)::json as sensors
								from instruments i
								left join cte_readings cte_r on i.instrument_id = cte_r.instrument_id
								where i.incinerator_id = ? and i.instrument_id = ? and i.deleted_at is null
								group by
									i.incinerator_id,
									i.instrument_id,
									i.instrument_name,
									i.instrument_code,
									cte_r.measure
				)
		select 
				gm.incinerator_id,
				gm.instrument_id,
				gm.instrument_name,
				gm.instrument_code,
				json_agg(
						json_build_object(
								'measure', gm.measure,
								'sensors', gm.sensors
						)::json
				)::json as measurement
		from grouped_measurement as gm
		group by
				gm.incinerator_id,
				gm.instrument_id,
				gm.instrument_name,
				gm.instrument_code
)`)
	return stmt
}
func getCteInstrumentV4() string {
	stmt := fmt.Sprintln(`cte_instruments as (
		with grouped_measurement as (
					select
							i.incinerator_id,
							i.instrument_id,
							i.instrument_name,
							i.instrument_code,
							cte_r.measure,
							json_agg(
															json_build_object(
																			'sensor_id', cte_r.sensor_id,                                
																			'label', cte_r.label,
																			'uom', cte_r.unit_of_measurement,
																			'sensor_values', cte_r.values
															)::json
							order by cte_r.label)::json as sensors
							from instruments i
							left join cte_readings cte_r on i.instrument_id = cte_r.instrument_id
							where i.incinerator_id = ? and i.instrument_id = ? and i.deleted_at is null
							group by
									i.incinerator_id,
									i.instrument_id,
									i.instrument_name,
									i.instrument_code,
									cte_r.measure
						)
		select 
						gm.incinerator_id,
						gm.instrument_id,
						gm.instrument_name,
						gm.instrument_code,
						json_agg(
										json_build_object(
														'measure', gm.measure,
														'sensors', gm.sensors
										)::json
						)::json as measurement
		from grouped_measurement as gm
		group by
						gm.incinerator_id,
						gm.instrument_id,
						gm.instrument_name,
						gm.instrument_code
	)`)
	return stmt
}

func getCteIncinerator() string {
	stmt := fmt.Sprintln(`cte_incinerators as (
		select
			i.destination_id,
			i.incinerator_id,
			i.incinerator_code,
			i.line,
			json_agg(
				json_build_object(
					'instrument_id', cte_i.instrument_id,
					'instrument_name', cte_i.instrument_name,
					'instrument_code', cte_i.instrument_code,
					'sensors', cte_i.sensors
				)::json
			)::json as instruments
		from incinerators i
		left join cte_instruments cte_i on i.incinerator_id = cte_i.incinerator_id
		where i.incinerator_id = ?
		group by
			i.destination_id,
			i.incinerator_id,
			i.incinerator_code,
			i.line
	)`)
	return stmt
}
func getCteIncineratorV2() string {
	stmt := fmt.Sprintln(`cte_incinerators as (
		select
				i.destination_id,
				i.incinerator_id,
				i.incinerator_code,
				i.line,
				json_agg(
						json_build_object(
								'instrument_id', cte_i.instrument_id,
								'instrument_name', cte_i.instrument_name,
								'instrument_code', cte_i.instrument_code,
								'measurement', cte_i.measurement
						)::json
				)::json as instruments
		from incinerators i
		left join cte_instruments cte_i on i.incinerator_id = cte_i.incinerator_id
		where i.incinerator_id = ?
		group by
				i.destination_id,
				i.incinerator_id,
				i.incinerator_code,
				i.line
)`)
	return stmt
}

// ________________________________________________________________________
func getCteReadingByDay() string {
	stmt := fmt.Sprintln(`
	daily_series as (select generate_series(?::timestamp, ?::timestamp, '1 day'::interval) as day)
	, derived_readings as (
		select 
			s.sensor_id, 
			s.label, 
			s.unit_of_measurement, 
			s.measure, 
			ds.day as read_at, 
			coalesce(round(avg(r.value::numeric), 3),0) as value
		from sensors s
		cross join daily_series ds
		left join readings r 
			on r.sensor_id = s.sensor_id 
			and ds.day = date_trunc('day', r.read_at) 
			and date(r.read_at) between ? and ?        
		group by 1,2,3,4,5
	)
	,cte_readings as (
		select
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement,
			json_agg(
					json_build_object(
						'read_at', r.read_at,
						'value', r.value
					)::json
			order by r.read_at)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement
)
	`)
	return stmt
}
func getCteReadingByWeek() string {
	stmt := fmt.Sprintln(`
	weekly_series as (
		with weekly as ( select generate_series(?::timestamp, ?::timestamp, '7 day'::interval) as weeks )
		select weeks, extract(week from weeks) as week_number from weekly
		) -- group by week num
, derived_readings as (
	select 
		s.sensor_id, 
		s.label, 
		s.unit_of_measurement, 
		s.measure, 
		concat('Week ', ws.week_number , ' - ', date(ws.weeks)) as read_at,
		coalesce(round(avg(r.value::numeric),3), 0) as value
	from sensors s
	cross join weekly_series ws
	left join readings r
		on r.sensor_id = s.sensor_id
			and ws.week_number = extract('week' from r.read_at)
			and date(r.read_at) between ? and ?
		group by 1,2,3,4,5
	)
	,cte_readings as (
		select
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement,
			json_agg(
					json_build_object(
						'read_at', r.read_at,
						'value', r.value
					)::json
			order by r.read_at)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement
)
	`)
	return stmt
}
func getCteReadingByMonth() string {
	stmt := fmt.Sprintln(`
		monthly_series as (
		with monthly as ( select generate_series(?::date, ?::date, '1 month'::interval) as months )
		select months, to_char(months, 'mon') as month from monthly
		) -- group by month num
, derived_readings as (
	select 
		s.sensor_id, 
		s.label, 
		s.unit_of_measurement, 
		s.measure, 
		ms.month as read_at,
		ms.months,
		coalesce(round(avg(r.value::numeric),3), 0) as value
	from sensors s
	cross join monthly_series ms
	left join readings r
		on r.sensor_id = s.sensor_id
			and extract('month' from ms.months) = extract('month' from r.read_at)
			and date(r.read_at) between ? and ?
		group by 1,2,3,4,5,6
	)
	,cte_readings as (
		select
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement,
			json_agg(
				json_build_object(
					'read_at', r.read_at,
					'value', r.value
				)::json
				order by r.months)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement
)
	`)
	return stmt
}
func getCteReadingByYear() string {
	stmt := fmt.Sprintln(`
	yearly_series as (
		with yearly as (select generate_series(?::date, ?::date, '1 year'::interval) as years)
			select years, extract(year from years) as year from yearly
		) -- group by year
, derived_readings as (
		select 
			s.sensor_id, 
			s.label, 
			s.unit_of_measurement, 
			s.measure, 
			ys.year as read_at,
			coalesce(round(avg(r.value::numeric),3), 0) as value
		from sensors s
		cross join yearly_series ys
		left join readings r
			on r.sensor_id = s.sensor_id
				and ys.year = extract('year' from r.read_at)
				and date(r.read_at) between ? and ?
			group by 1,2,3,4,5
	)
	,cte_readings as (
		select
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement,
			json_agg(
					json_build_object(
						'read_at', r.read_at,
						'value', r.value
					)::json
			order by r.read_at)::json as values
		from derived_readings r
		left join sensors s on s.sensor_id = r.sensor_id
		group by
			s.instrument_id,
			s.measure,
			s.label,
			s.sensor_id,
			s.unit_of_measurement
)
	`)
	return stmt
}
