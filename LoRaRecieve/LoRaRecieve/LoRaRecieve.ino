#include <SPI.h>
#include <LoRa.h>
#include <SoftwareSerial.h>

SoftwareSerial esp(6, 7); // RX, TX



String rssi = "RSSI --";
String packSize = "--";
String packet ;

const int csPin = 10;          // LoRa radio chip select
const int resetPin = 5;       // LoRa radio reset
const int irqPin = 2;         // change for your board; must be a hardware interrupt pin

void cbk(int packetSize) {
  packet ="";
  packSize = String(packetSize,DEC);
  for (int i = 0; i < packetSize; i++) { packet += (char) LoRa.read(); }
  rssi = "RSSI " + String(LoRa.packetRssi(), DEC) ;
  int index1 = packet.indexOf(',');
  int index2 = packet.indexOf(',', index1+1);
  String node = packet.substring(0,index1);
  String temp = packet.substring(index1+1,index2);
  String co = packet.substring(index2+1);
  esp.print(packet);
  Serial.println(packet + "  " + rssi);
  Serial.println("Node = " + node + "\tTemp = " + temp + "\tCO2 = " + co);
}
 
void setup() {
  esp.begin(9600);
  Serial.begin(9600);
  while (!Serial);
  Serial.println("LoRa Receiver");
  LoRa.setPins(csPin, resetPin, irqPin); // set CS, reset, IRQ pin
  if (!LoRa.begin(915E6)) {
    Serial.println("Starting LoRa failed!");
    while (1);
  }
  LoRa.setSyncWord(0x11);           // ranges from 0-0xFF, default 0x34, see API docs
  LoRa.setSpreadingFactor(12);
  Serial.println("LoRa Start OK");
}

void loop() {
  // try to parse packet

  int packetSize = LoRa.parsePacket();
  if (packetSize) {
    cbk(packetSize);
  }
  delay(10);
}
