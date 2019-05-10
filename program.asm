// Add Skip size of image
ADDI V2 V2 #4

#repeat 1000
LOAD V5 V2
// Editar la imagen
// ADDI V5 V0 #100
// ADD Encriptation
// ADD255 V5 V5 #50
// ADD255 V5 V5 #-50
// XOR Encriptation
// XOR255 V5 V5 #65
// XOR255 V5 V5 #65
// SHUFFLE V5 V5
// UNSHUFFLE V5 V5
// SHUFFLE255 V5 V5
// UNSHUFFLE255 V5 V5
FLIP V5 V5
// FLIP V5 V5
STORE V2 V5
ADDI V2 V2 #8
#endrepeat
