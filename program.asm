// Add Skip size of image
ADDI V2 V2 #4

#repeat 9256
LOAD V5 V2
// Editar la imagen
// ADDI V5 V0 #100
// ADD Encriptation
// ADD255 V5 V5 #155
// ADD255 V5 V5 #-155
// XOR Encriptation
// XOR255 V5 V5 #170
// XOR255 V5 V5 #170
// SHUFFLE V5 V5
// UNSHUFFLE V5 V5
SHUFFLE255 V5 V5
UNSHUFFLE255 V5 V5
// FLIP V5 V5
// FLIP V5 V5
// UNSHUFFLE255 V5 V5
// RL V5 V5 #55
// RR V5 V5 #55
STORE V2 V5
ADDI V2 V2 #8
#endrepeat
