/*
  LoRa Duplex communication with Sync Word

  Sends a message every half second, and polls continually
  for new incoming messages. Sets the LoRa radio's Sync Word.

  Spreading factor is basically the radio's network ID. Radios with different
  Sync Words will not receive each other's transmissions. This is one way you
  can filter out radios you want to ignore, without making an addressing scheme.

  See the Semtech datasheet, http://www.semtech.com/images/datasheet/sx1276.pdf
  for more on Sync Word.

  created 28 April 2017
  by Tom Igoe
*/
#include <SPI.h>              // include libraries
#include <LoRa.h>
const int csPin = 10;          // LoRa radio chip select
const int resetPin = 5;       // LoRa radio reset
const int irqPin = 2;         // change for your board; must be a hardware interrupt pin

byte msgCount = 0;            // count of outgoing messages
int interval = 2000;          // interval between sends
long lastSendTime = 0;        // time of last packet send

String node = "1";

void setup() {
  Serial.begin(9600);                   // initialize serial
  while (!Serial);
  
  Serial.println("LoRa Duplex - Set sync word");

  // override the default CS, reset, and IRQ pins (optional)
  LoRa.setPins(csPin, resetPin, irqPin);// set CS, reset, IRQ pin

  if (!LoRa.begin(915E6)) {             // initialize ratio at 915 MHz
    Serial.println("LoRa init failed. Check your connections.");
    while (true);                       // if failed, do nothing
  }

  LoRa.setSyncWord(0x11);           // ranges from 0-0xFF, default 0x34, see API docs
  LoRa.setSpreadingFactor(12);
  Serial.println("LoRa init succeeded.");
}

void loop() {
  LoRa.beginPacket();
  LoRa.print(node + "," + fakesensortemp() + "," + fakesensorco());
  LoRa.endPacket();
  Serial.println("send OK");
  delay(2000);
}


void sendMessage(String outgoing) {
  LoRa.beginPacket();                   // start packet
  LoRa.print(outgoing);                 // add payload
  LoRa.endPacket();                     // finish packet and send it
}

String fakesensortemp() {
    int ran = random(250,400);
    return String(ran/10.0);
}

String fakesensorco() {
    int ran = random(100,1000);
    return String(ran);
}
