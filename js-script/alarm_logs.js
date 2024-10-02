require("dotenv").config({
  path: "../.env",
});
const mqtt = require("mqtt");

const protocol = "mqtt";
const host = process.env.MQ_HOST;
const port = process.env.MQ_PORT;
const username = process.env.MQ_USER;
const password = process.env.MQ_PASS;

console.log(password);

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


const incineratorList = [
  {
    destination_code: "03",
    incinerator_code: "line a",
      instrument_name: "Rotary Klin",
      instrument_code: "RK",  
    },
    {
      destination_code: "03",
      incinerator_code: "line a",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",  
    },
    {
      destination_code: "03",
      incinerator_code: "line a",
      instrument_name: "Water Quencher",
      instrument_code: "WQ",  
    },
    {
      destination_code: "03",
      incinerator_code: "line b",
      instrument_name: "Rotary Klin",
      instrument_code: "RK",  
    },
    {
      destination_code: "03",
      incinerator_code: "line b",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",  
    },
    {
      destination_code: "03",
      incinerator_code: "line b",
      instrument_name: "Water Quencher",
      instrument_code: "WQ",  
    },
    {
      destination_code: "04",
      incinerator_code: "line a",
      instrument_name: "Rotary Klin",
      instrument_code: "RK",  
    },
    {
      destination_code: "04",
      incinerator_code: "line a",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",  
    },
    {
      destination_code: "04",
      incinerator_code: "line a",
      instrument_name: "Water Quencher",
      instrument_code: "WQ",  
    },
    {
      destination_code: "04",
      incinerator_code: "line b",
      instrument_name: "Rotary Klin",
      instrument_code: "RK",  
    },
    {
      destination_code: "04",
      incinerator_code: "line b",
      instrument_name: "Secondary Combustion Chamber",
      instrument_code: "SCC",  
    },
    {
      destination_code: "04",
      incinerator_code: "line b",
      instrument_name: "Gas Cooler",
      instrument_code: "GC",  
    },
  ];
  
  const topic = "afes/iot/scada/alarm_logs";
  client.on("connect", () => {
    console.log("Connected to MQTT broker");

    const startDate = new Date("2023-08-01"); // Replace with your desired start date
    const endDate = new Date("2023-08-30"); // Replace with your desired end date

    incineratorList.forEach(incinerator => {
      let currentDate = new Date(startDate)

      while (currentDate <= endDate){
        const message = {
          destination_code: incinerator.destination_code,
          incinerator_code: incinerator.incinerator_code,
          instrument_name: incinerator.instrument_name,
          instrument_code: incinerator.instrument_code,
          time: formatDateTime(currentDate),
          priority: "1",
          state: "UNACK",
          node: "DHES-SCADA",
          group: "$System",
          tag_name: "TC2_RAW",
          description: incinerator.instrument_name,
          type: "LO",
          limit: "850",
          current_value: "0",
          alarm_duration: "000 00:00:00.010",
          operator: "None",
          un_ack_duration: "000 00:00:00.100",
        };
      
        // Loop
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

        currentDate.setDate(currentDate.getDate() + 1)
      }

    })

  client.end();
});
client.on("error", (error) => {
  console.error("MQTT connection error:", error);
});

// "6/29/2021 4:22:56"
function formatDateTime(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");
  const milliseconds = String(date.getMilliseconds()).padStart(3, "0");

  return `${month}/${day}/${year} ${hours}:${minutes}:${seconds}`
}