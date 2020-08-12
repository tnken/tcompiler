#define LED_PIN 13

void setup() {
  pinMode(LED_PIN, OUTPUT);
  Serial.begin(9600);
}

void loop() {
  if(Serial.available() > 0 ){
    int incomingByte = Serial.read();
    switch(incomingByte){
      case '1':
        digitalWrite(LED_PIN, HIGH);
        break;
      case '0':
        digitalWrite(LED_PIN, LOW);
        break;
    }
  }
}
