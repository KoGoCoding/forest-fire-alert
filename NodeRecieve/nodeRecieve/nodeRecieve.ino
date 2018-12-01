#include <SoftwareSerial.h>

#include <LiquidCrystal_I2C.h>

#include <Wire.h>
//#include <WiFi.h>
//#include <WiFiClientSecure.h>
#include <SPI.h>
#include <ESP8266WiFi.h>
#include <WiFiClient.h> 
#include <ESP8266WebServer.h>
#include <ESP8266HTTPClient.h>


#define WIFI_SSID "AndroidAPz"
#define WIFI_PASSWORD "Ko157953"

#define relay D3
#define sw D4

LiquidCrystal_I2C lcd(0x27, 20, 4);

SoftwareSerial arduino(D6, D7); // RX, TX

int status_fire = 0;

int state = 0;



void setup() {
  WiFi.mode(WIFI_STA);
  // connect to wifi.
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  Serial.print("connecting");

  while (WiFi.status() != WL_CONNECTED) {
    Serial.print(".");
    delay(500);
  }
  Serial.println();
  Serial.print("connected: ");
  Serial.println(WiFi.localIP());
  arduino.begin(9600);
  Serial.begin(9600);
  pinMode(sw, INPUT_PULLUP);
  pinMode(relay, OUTPUT);

  Wire.begin(D2, D1);

  lcd.begin();
  
  lcd.home();
  
  lcd.print("Forest Fire Aleart");
  lcd.setCursor(0,1);
}

String data = "";
String temp = "";
void loop() {
  HTTPClient http;    //Declare object of class HTTPClient
  if(arduino.available()){
    data = arduino.readString();
    Serial.println(data);
    int index1 = data.indexOf(',');
    int index2 = data.indexOf(',', index1+1);
    String node = data.substring(0,index1);
    temp = data.substring(index1+1,index2);
    String co = data.substring(index2+1);
    String co2_status = "";

    if(co.toInt() < 1800) {
      co2_status = "LOW";
    } else if (co.toInt() < 2300) {
      co2_status = "WARNING";
    } else if (co.toInt() > 2300) {
      co2_status = "DANGER";
    }
    
    lcd.setCursor(0,1);
    lcd.print("Node1 Temp = "+temp+"*C CO2 = "+co2_status+"       ");


    
    if(temp.toInt() > 35 && status_fire != 4) {
      status_fire = 2;
    } else if (temp.toInt() > 30 && status_fire == 0) {
      status_fire = 1;
    } else if (temp.toInt() < 28) {
      status_fire = 0;
    }

//    //POST request
//      String postData = "co2=" + co + "&temp=" + temp+ "&nodeID=1" ;
//     http.begin("http://192.168.43.103/savedata");              //Specify request destination
//    http.addHeader("Content-Type", "application/x-www-form-urlencoded");    //Specify content-type header
//   
//    int httpCode = http.POST(postData);   //Send the request
//    String payload = http.getString();    //Get the response payload
//    http.end();
    
  }

  int val_sw = digitalRead(sw);
  val_sw = !val_sw;
  Serial.println(val_sw);

  if(val_sw == 1) {
    if(status_fire == 1) {
      status_fire = 3;
    } else if (status_fire == 2) {
      status_fire = 4;
    }
  }

  
  if(status_fire == 1) {
    //เข้าไลน์
  } else if(status_fire == 2 && state != 1) {
    digitalWrite(D3, HIGH);
    state = 1;
     http.begin("http://192.168.43.103:8080/alert");              //Specify request destination
    http.addHeader("Content-Type", "application/x-www-form-urlencoded");    //Specify content-type header
   
    int httpCode = http.POST("");   //Send the request
//    String payload = http.getString();    //Get the response payload
    http.end();
  } else if(status_fire == 4) {
    digitalWrite(D3, LOW);
  } else if(status_fire == 0) {
    digitalWrite(D3, LOW);
    state = 0;
  }
  Serial.println(status_fire);
}


// http.begin("http://192.168.43.103:8080/alert");              //Specify request destination
//    http.addHeader("Content-Type", "application/x-www-form-urlencoded");    //Specify content-type header
//   
//    int httpCode = http.POST("");   //Send the request
////    String payload = http.getString();    //Get the response payload
//    http.end();
