# blogvm - A simple VM project to go with a blog series
I'm currently writing a blog series to go with this project, the
first of which is available [here](https://www.stephengream.com/writing-a-vm-part-one)

## Instruction structure
Our computer is a 32-bit computer, meaning our control unit executes instructions that are
32 bits in length. We will use hexadecimal here to simplify reading.

![What an instruction looks like](./docs/Instruction.png)

We break our instructions up into four parts:
* **The Opcode** - One byte for the instruction to execute, for example an add or multiply operation
* **Input 1 (`I1`)** - 4 bits for the first input register, we will elaborate on how to address these later
* **Input 2/Destination (`I2`/`D`)** - 4 bits for the second input register, which will double as a destination value.
* **Immediate data** - 2 bytes for any immediate data. More on this later

## Instruction Reference
| Hex Value | Mnemonic | Description                                          |
|-----------|----------|------------------------------------------------------|
| 0x00      | HALT     | Halt the machine                                     |
| 0x01      | READ     | Read from the memory address in I1                   |
| 0x02      | WRITE    | Write to the memory address in D                     |
| 0x03      | COPY     | Copy from register I1 to D                           |
| 0x04      | ADD      | Add I1 to I2 and store in D                          |
| 0x05      | SUB      | Subtract I2 from I1 and store the result in D        |
| 0x06      | MUL      | Multiply I1 by I2 and store the result in D          |
| 0x07      | DIV      | Divide I1 by I2 and store the result in D            |
| 0x08      | STAT     | Get the status at bit I1 and store the result in D   |
| 0x09      | SET      | Set the status at but I1 to I2                       |


## Addressing registers
Each register is addressed with 4 bits of data, meaning we can potentially refer to 16 different
registers. We will use the value 0xF to refer to the **immediate data** in an instruction

| Hex Value | Mnemonic | Description          | Initial Value |
|-----------|----------|----------------------|---------------|
| 0x0       | R0       | First ALU register   | 0x00          |
| 0x1       | R1       | Second ALU register  | 0x00          |
| 0x2       | R2       | Third ALU register   | 0x00          |
| 0x3       | R4       | Fourth ALU register  | 0x00          |
| 0xC       | SR       | Status register      | 0x00          |
| 0xD       | PC       | Program counter      | 0x100         |
| 0xE       | IR       | Instruction register | 0x00          |
| 0xF       | #{n}     | Immediate data       | N/A           |
