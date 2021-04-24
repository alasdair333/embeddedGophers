package main

import (
	"unsafe"
)

const (
	mmioBaseAddr   uintptr = 0x3F000000
	AUX_ENABLE     uintptr = (mmioBaseAddr + 0x00215004)
	AUX_MU_IO      uintptr = (mmioBaseAddr + 0x00215040)
	AUX_MU_IER     uintptr = (mmioBaseAddr + 0x00215044)
	AUX_MU_IIR     uintptr = (mmioBaseAddr + 0x00215048)
	AUX_MU_LCR     uintptr = (mmioBaseAddr + 0x0021504C)
	AUX_MU_MCR     uintptr = (mmioBaseAddr + 0x00215050)
	AUX_MU_LSR     uintptr = (mmioBaseAddr + 0x00215054)
	AUX_MU_MSR     uintptr = (mmioBaseAddr + 0x00215058)
	AUX_MU_SCRATCH uintptr = (mmioBaseAddr + 0x0021505C)
	AUX_MU_CNTL    uintptr = (mmioBaseAddr + 0x00215060)
	AUX_MU_STAT    uintptr = (mmioBaseAddr + 0x00215064)
	AUX_MU_BAUD    uintptr = (mmioBaseAddr + 0x00215068)

	GPFSEL1   uintptr = (mmioBaseAddr + 0x00200004)
	GPPUD     uintptr = (mmioBaseAddr + 0x00200094)
	GPPUDCLK0 uintptr = (mmioBaseAddr + 0x00200098)
)

// func noop() {

// }

func uart_init() int {
	*(*uintptr)(unsafe.Pointer(AUX_ENABLE)) = 1 // enable UART1, AUX mini uart
	*(*uintptr)(unsafe.Pointer(AUX_MU_CNTL)) = 0
	*(*uintptr)(unsafe.Pointer(AUX_MU_LCR)) = 3 // 8 bits
	*(*uintptr)(unsafe.Pointer(AUX_MU_MCR)) = 0
	*(*uintptr)(unsafe.Pointer(AUX_MU_IER)) = 0
	*(*uintptr)(unsafe.Pointer(AUX_MU_IIR)) = 0xc6 // disable interrupts
	*(*uintptr)(unsafe.Pointer(AUX_MU_BAUD)) = 270 // 115200 baud
	/* map UART1 to GPIO pins */
	//r := *(*uintptr)(unsafe.Pointer(GPFSEL1))
	//r &= (7 << 12) | (7 << 15) // gpio14, gpio15
	//r |= (2 << 12) | (2 << 15)

	//*(*uintptr)(unsafe.Pointer(GPFSEL1)) = r
	*(*uintptr)(unsafe.Pointer(GPPUD)) = 0 // enable pins 14 and 15

	*(*uintptr)(unsafe.Pointer(GPPUDCLK0)) = (1 << 14) | (1 << 15)

	*(*uintptr)(unsafe.Pointer(GPPUDCLK0)) = 0   // flush GPIO setup
	*(*uintptr)(unsafe.Pointer(AUX_MU_CNTL)) = 3 // enable Tx, Rx
	return 0
}

func uart_send(c byte) {
	*(*byte)(unsafe.Pointer(AUX_MU_IO)) = c
}

func uart_puts(s []byte) {

	for _, c := range s {
		uart_send(c)
	}
}

func main() {

	uart_init()

	uart_send('H')
	uart_send('E')
	uart_send('L')
	uart_send('L')
	uart_send('O')
	uart_send(' ')
	uart_send('G')
	uart_send('O')
	uart_send('!')

	for {

	}
}
