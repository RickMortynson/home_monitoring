#include <Arduino.h>
#include <ArduinoOTA.h>
#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClientSecure.h>

const char *WIFI_SSID = "ssid";
const char *WIFI_PASSWORD = "pass";

const String host = "https://heartbeat_Service.heartbeat.sh";
const String url = "/beat/ping";
const int port = 443;

WiFiClientSecure wifiClient;
HTTPClient httpClient;

void setup() {
  Serial.begin(115200);
  Serial.println("Booting");

  // WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  while (WiFi.waitForConnectResult() != WL_CONNECTED) {
    Serial.println("Connection Failed! Rebooting...");
    delay(2000);
    ESP.restart();
  }

  // 
  // OTA START
  // 

  ArduinoOTA.onStart([]() {
    Serial.println("Start OTA \n");
  });
  ArduinoOTA.onEnd([]() {
    Serial.println("\nEnd OTA");
  });
  ArduinoOTA.onProgress([](unsigned int progress, unsigned int total) {
    Serial.printf("Progress: %u%%\r", (progress / (total / 100)));
  });
  ArduinoOTA.onError([](ota_error_t error) {
    Serial.printf("Error[%u]: ", error);
    if (error == OTA_AUTH_ERROR) Serial.println("Auth Failed");
    else if (error == OTA_BEGIN_ERROR) Serial.println("Begin Failed");
    else if (error == OTA_CONNECT_ERROR) Serial.println("Connect Failed");
    else if (error == OTA_RECEIVE_ERROR) Serial.println("Receive Failed");
    else if (error == OTA_END_ERROR) Serial.println("End Failed");
  });
  ArduinoOTA.begin();
  Serial.println("Ready");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

  // 
  // OTA END
  // 
  
  wifiClient.setInsecure();
  wifiClient.setNoDelay(true);
  wifiClient.connect(host, port);
}



void loop() {
    ArduinoOTA.handle();

    httpClient.begin(wifiClient, host + url);
  httpClient.setTimeout(1000);

    Serial.println("POST...");
    int httpCode = httpClient.POST("/");
    Serial.println(httpCode);
    httpClient.end();
    if (httpCode == 200) {
      // restart immediately if got any error code
      delay(3500);
    } else {
      delay(1000);
    }
}