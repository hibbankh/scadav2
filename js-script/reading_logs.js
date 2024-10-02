require("dotenv").config({
  path: "../.env"
});
const mqtt = require("mqtt");

const protocol = "mqtt";
const host = process.env.MQ_HOST;
const port = process.env.MQ_PORT;
const username = process.env.MQ_USER;
const password = process.env.MQ_PASS;

const clientId = `mqtt_${Math.random().toString(16).slice(3)}`;
const connectUrl = `${protocol}://${host}:${port}`;

const client = mqtt.connect(connectUrl, {
  clientId,
  clean: true,
  connectTimeout: 4000,
  username,
  password,
  reconnectPeriod: 1000,
});

const CH_A = {
  destination_code: "03",
  incinerator_code: "line a",
  instrument: [
    {
      instrument_name: "Rotary Klin",
      instrument_code: "RK",
      sensor: [
        { label: "TT-301A", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-301B", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-401A", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-401B", unit_of_measurement: "c", measure: "temperature" },
        { label: "PT-401", unit_of_measurement: "mba", measure: "pressure" },
        { label: "VFD", unit_of_measurement: "hz", measure: "frequency" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",
      sensor: [
        { label: "PT-501", unit_of_measurement: "mba", measure: "pressure" },
        { label: "TT-503", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-501A", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-501B", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-505", unit_of_measurement: "c", measure: "temperature" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Water Quencher",
      instrument_code: "WQ",
      sensor: [
        { label: "DP-501", unit_of_measurement: "mba", measure: "pressure" },
        { label: "TT-502", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-502B", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-504", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-501A", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-501B", unit_of_measurement: "c", measure: "temperature" },
        { label: "DP-601", unit_of_measurement: "mba", measure: "pressure" },
        { label: "DP-701", unit_of_measurement: "mba", measure: "pressure" },
        { label: "TT-702", unit_of_measurement: "c", measure: "temperature" },
        { label: "TT-701", unit_of_measurement: "c", measure: "temperature" },
        { label: "VFD", unit_of_measurement: "hz", measure: "frequency" },
      ],
    },
  ],
};
const CH_B = {
  destination_code: "03",
  incinerator_code: "line b",
  instrument: [
    {
      instrument_name: "Rotary Klin",
      instrument_code: "RK",
      sensor: [
        { label: "BBTE05", unit_of_measurement: "c", measure: "temperature" },
        { label: "BBTE01", unit_of_measurement: "c", measure: "temperature" },
        { label: "BBTE02", unit_of_measurement: "c", measure: "temperature" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",
      sensor: [
        { label: "BBTE04", unit_of_measurement: "c", measure: "temperature" },
        { label: "BBTE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "CBPT01", unit_of_measurement: "mba", measure: "pressure" },
        { label: "CBTE02", unit_of_measurement: "c", measure: "temperature" },
        { label: "CBTE01", unit_of_measurement: "c", measure: "temperature" },
        { label: "CBPT02", unit_of_measurement: "mba", measure: "pressure" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Water Quencher",
      instrument_code: "WQ",
      sensor: [
        { label: "CBTE04", unit_of_measurement: "c", measure: "temperature" },
        { label: "CBTE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "DBPT01", unit_of_measurement: "mba", measure: "pressure" },
        { label: "DBTE01", unit_of_measurement: "c", measure: "temperature" },
        { label: "DBPT02", unit_of_measurement: "mba", measure: "pressure" },
        { label: "DBTE02", unit_of_measurement: "c", measure: "temperature" },
        { label: "DBPT03", unit_of_measurement: "mba", measure: "pressure" },
        { label: "DBTE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "DBFN01", unit_of_measurement: "hz", measure: "frequency" },
      ],
    },
  ],
};
const PP_A = {
  destination_code: "04",
  incinerator_code: "line a",
  instrument: [
    {
      instrument_name: "Rotary Klin",
      instrument_code: "RK",
      sensor: [
        { label: "BAPT01", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "BAFN01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "BARD01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "BATE01", unit_of_measurement: "c", measure: "temperature" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",
      sensor: [
        { label: "BATE02", unit_of_measurement: "c", measure: "temperature" },
        { label: "BATE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "BATE04", unit_of_measurement: "c", measure: "temperature" },
        { label: "BAFN03", unit_of_measurement: "hz", measure: "frequency" },
        { label: "CATE01", unit_of_measurement: "c", measure: "temperature" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Water Quencher",
      instrument_code: "WQ",
      sensor: [
        { label: "CATE05", unit_of_measurement: "c", measure: "temperature" },
        { label: "CATE04", unit_of_measurement: "c", measure: "temperature" },
        { label: "CATE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCDC02", unit_of_measurement: "hz", measure: "frequency" },
        { label: "DCDC01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "DCPT02", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCTE01", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCTE02", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCPT03", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCTE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCFN01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "DCPT04", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCTE04", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCFN02", unit_of_measurement: "hz", measure: "frequency" },
        { label: "CAPT02", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "CAPT01", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "CAPU02", unit_of_measurement: "hz", measure: "frequency" },
        { label: "CAPU01", unit_of_measurement: "hz", measure: "frequency" },
      ],
    },
  ],
};

const PP_B = {
  destination_code: "04",
  incinerator_code: "line b",
  instrument: [
    {
      instrument_name: "Rotary Klin",
      instrument_code: "RK",
      sensor: [
        { label: "BBPT01", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "BBFN01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "BBTE01", unit_of_measurement: "c", measure: "temperature" },
        { label: "BBRD01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "BBTE02", unit_of_measurement: "c", measure: "temperature" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",
      sensor: [
        { label: "BBTE03", unit_of_measurement: "c", measure: "temperature" },
        { label: "BBTE04", unit_of_measurement: "c", measure: "temperature" },
        { label: "BBFN02", unit_of_measurement: "hz", measure: "frequency" },
        { label: "CBPT02", unit_of_measurement: "mBar", measure: "pressure" },
      ],
    },
    {
      incinerator_code: "line a",
      instrument_name: "Gas Cooler",
      instrument_code: "GC",
      sensor: [
        { label: "CBTE02", unit_of_measurement: "c", measure: "temperature" },
        { label: "CBPT01", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "CBFN01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "CBFN02", unit_of_measurement: "hz", measure: "frequency" },
        { label: "CBTE01", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCPT05", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCPT06", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCPT07", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCPT08", unit_of_measurement: "mBar", measure: "pressure" },
        { label: "DCTE06", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCTE07", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCTE08", unit_of_measurement: "c", measure: "temperature" },
        { label: "DCFN01", unit_of_measurement: "hz", measure: "frequency" },
        { label: "DCFN02", unit_of_measurement: "hz", measure: "frequency" },
      ],
    },
  ],
};


const topic = "afes/iot/scada/reading_logs";
client.on("connect", () => {
  console.log("Connected to MQTT broker");

  const incinerators = [CH_A, CH_B, PP_A, PP_B];
  incinerators.forEach((incinerator) => {
    const startDate = new Date("2023-08-21"); // Replace with your desired start date
    const endDate = new Date("2023-08-30"); // Replace with your desired end date

    let currentDate = new Date(startDate)
    while (currentDate <= endDate){

        // make a copy
        const copy = incinerator;
        copy.instrument.forEach((instrument) => {
          instrument.sensor.forEach((sensor) => {
            // Generate random values for hours, minutes, seconds, and milliseconds
            const randomHours = Math.floor(Math.random() * 24);
            const randomMinutes = Math.floor(Math.random() * 60);
            const randomSeconds = Math.floor(Math.random() * 60);
            const randomMilliseconds = Math.floor(Math.random() * 1000); // Microseconds are not directly supported in JavaScript
    
            // Set the random time components
            currentDate.setHours(
              randomHours,
              randomMinutes,
              randomSeconds,
              randomMilliseconds
            );
            // generate random datetime
            if (sensor.measure === "temperature") {
              sensor["value"] = Math.floor(Math.random() * (5555 + 1)).toString();
            } else if (sensor.measure === "pressure") {
              sensor["value"] = (Math.floor(Math.random() * (100 + 1 - -10)) - -10).toString();
            } else if (sensor.measure === "frequency") {
              sensor["value"] = Math.floor(Math.random() * (50 + 1)).toString();
            }
            sensor["read_at"] = formatDateTime(currentDate);
          });
          const message = {
            destination_code: incinerator.destination_code,
            incinerator_code: incinerator.incinerator_code,
            instrument_name: instrument.instrument_name,
            instrument_code: instrument.instrument_code,
            sensor: instrument.sensor,
          };
          client.publish(
            topic,
            JSON.stringify(message),
            { qos: 0, retain: false },
            (error) => {
              if (error) {
                console.error("Error publishing message:", error);
              } else {
                console.log("Message published successfully");
              }
            }
          );
        });

        
        currentDate.setDate(currentDate.getDate() + 1)
    }

  });
  client.end();
});
client.on("error", (error) => {
  console.error("MQTT connection error:", error);
});

function formatDateTime(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");
  const milliseconds = String(date.getMilliseconds()).padStart(3, "0");

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}`;
}
