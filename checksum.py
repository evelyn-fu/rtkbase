cls = 0x06
id = 0x01
msg_len = 0x14
payload = [0] * 20
payload[12] = 0x23
# payload[9] = 0x96
payload[3] = 0x8

checksumA = 0
checksumB = 0

checksumA += cls
checksumB += checksumA

checksumA += id
checksumB += checksumA

checksumA += (msg_len & 0xFF)
checksumB += checksumA

checksumA += (msg_len >> 8)
checksumB += checksumA

for i in range(msg_len):
    checksumA += payload[i]
    checksumB += checksumA

print("Checksum A: ", hex(checksumA))
print("Checksum B: ", hex(checksumB))

synch1 = 0xB5
synch2 = 0x62
lsb = 0x14
msb = 0x00

message = [synch1, synch2, cls, id, lsb, msb] + payload + [checksumA % 256, checksumB % 256]

print(message)