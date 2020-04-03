package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	CSVpath := "./FILES-CSV/"
	DATPath := "./FILES-DAT/"

	files, err := ioutil.ReadDir(CSVpath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == ".csv" {
				recordFile, err := os.Open(CSVpath + f.Name())
				if err != nil {
					fmt.Println("An error encountered ::", err)
				}
				reader := csv.NewReader(recordFile)
				reader.Comma = ';'
				// 3. Read all the records
				records, _ := reader.ReadAll()
				if len(records[0]) == 15 {

					stringToSave := ""
					rowHeader := records[1]
					var total float64 = 0
					sTotal := ""

					stringToSave += datHeader()
					records = records[1:]
					for _, record := range records {
						tmpValue, _ := strconv.ParseFloat(strings.Replace(record[5], ",", ".", -1), 64)
						total += tmpValue
						stringToSave += datBodyLine(record, rowHeader)
					}
					tmpRowsLines := strconv.Itoa((len(records) + 2))
					sTotal = strconv.FormatFloat(total, 'f', 0, 64) //le saco los decimales
					sTotal += "00"
					stringToSave += datFooter(tmpRowsLines, sTotal)

					//fmt.Println(records)
					//fmt.Println(datHeader)
					fmt.Println(stringToSave)
					saveFile(stringToSave, DATPath)

				} else {
					log.Fatal("El archivo CSV no tiene las 15 columnas necesarias")
				}
			}
		}
	}
}

func datHeader() string {
	//HLAVOZDELIN2020031723553500002071
	var header string
	t := time.Now()
	header = "H"
	header += "LAVOZDELIN" //codigo provisto SAP
	//header += t.Year() + t.Month() + t.Day() //Formato AAAAMMDD
	header += t.Format("20060102") //Formato AAAAMMDD
	header += t.Format("150405")   //Formato  HHMMSS
	header += "0000"               //Nro. identificador del archivo generado
	header += "2071"               //Nro. identificador del archivo recibido
	header += "\r\n"
	return header
}

func datFooter(rows string, total string) string {
	//HLAVOZDELIN2020031723553500002071
	var footer string
	footer += "T"                     //Identificación del registro
	footer += padLeft(rows, "0", 5)   // Cantidad de registros informados en el archivo  (incluyendo la Cabecera y este registro de Totales)
	footer += padLeft(total, "0", 15) // Sumatoria de los Importes al 1er. Vencimiento  (ver Registros de Deudas). Equivale al Total a Cobrar en los archivos que la Empresa le envía al Banco.
	footer += "\r\n"
	return footer
}

func datBodyLine(CVSRow []string, rowHeader []string) string {
	var tmpLine string
	layout := "02/01/2006"
	t, err := time.Parse(layout, CVSRow[1])
	check(err)
	var monedaDeuda string = "80"
	if CVSRow[4] != "ARS" {
		monedaDeuda = "02"
	}
	var cuil string = strings.Replace(CVSRow[11], "-", "", -1)
	cuil = strings.Replace(cuil, "$ ", "", -1)
	var importe string = strings.Replace(CVSRow[5], ",", "", -1)

	tmpLine = "D"                                    //Identificación del registro
	tmpLine += padLeft("42", "0", 4)                 //Código provisto por el Banco
	tmpLine += padLeft(" ", " ", 18)                 // 18c Código provisto por el Banco
	tmpLine += padLeft("0", "0", 18)                 // 18c Nro. de Factura
	tmpLine += padLeft("80", " ", 2)                 // 2c Tipo de Documento del Deudor
	tmpLine += padLeft(CVSRow[9], " ", 13)           // 13c Nro. de CUIT/CUIL o D.N.I. del Deudor
	tmpLine += padLeft(t.Format("20060102"), " ", 8) // 8c Fecha de origen de la Deuda (Puede informarse fecha inicio del convenio) Formato AAAAMMDD
	tmpLine += padLeft("20990908", " ", 8)           // 8c Fecha de caducidad de la Deuda (Fecha tope, puede ser año 2099) Formato AAAAMMDD
	tmpLine += padLeft("", "9", 15)                  // 15c (13,2 dec) Importe a pagar hasta el 1er. (o único) Vencimiento indicado
	tmpLine += padLeft(" ", " ", 8)                  // 8c Fecha del 1er. (o único) Vencimiento del pago ( Es Mandatario en caso de informar importe en campo anterior)
	tmpLine += padLeft(" ", "0", 15)                 // 15c (13,2 dec) Importe a pagar después del 1er. Vencimiento,  y hasta el 2do. Vto.
	tmpLine += padLeft(" ", " ", 8)                  // 8c Fecha del 2do. Vencimiento ( Es Mandatario en caso de informar importe en campo anterior)
	tmpLine += padLeft(" ", "0", 15)                 // 15c (13,2 dec) Importe a pagar después del 2do. Vencimiento, y hasta el 3er. Vto.
	tmpLine += padLeft(" ", " ", 8)                  // 8c Fecha del 3er. Vencimiento ( Es Mandatario en caso de informar importe en campo anterior)
	tmpLine += padRight(CVSRow[10], " ", 30)         // Apellido y nombre  o  Razón Social  del Deudor.
	tmpLine += padLeft(" ", " ", 30)                 // Código/Número que identifica a esta Deuda
	tmpLine += padLeft(" ", " ", 6)                  // Código o Número del  Area,  División,  Sucursal  o  Centro de Costo de la Empresa en el que se origina  la Deuda
	tmpLine += padLeft(" ", " ", 30)                 // Nombre del  Area,  División,  Sucursal  o  Centro de Costo  originante
	tmpLine += padLeft(monedaDeuda, " ", 2)          // Moneda de la Deuda
	tmpLine += padLeft(cuil, " ", 11)                // Nro. de CUIT de la Empresa,  División  o  Subsidiaria  en la que se origina la Deuda
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 38)                 // Hasta  10  lineas de texto  que aparecerán en el ticket de caja que entrega el Banco al momento del pago
	tmpLine += padLeft(" ", " ", 8)                  // Fecha de acreditación del pago recibido.  Es la fecha de cobranza para pagos en Efectivo o Débito en Cuenta, o la fecha en que se acreditó el cheque/ Importante: Si el Cheque se encuentra Rechazado, sus campos se informaran con Ceros y en Blanco cuando el cheque es depositado
	tmpLine += padLeft(" ", " ", 3)                  // Sucursal de cobro
	tmpLine += padLeft(" ", " ", 1)                  // Forma/Medio de pago: Es un código numérico de 1 dígito que indica si el pago se hizo en Efectivo, Débito en Cta., Cheques de este Banco, Cheques de 24 Hs., Cheques de Pago Diferido, etc.
	tmpLine += padLeft(" ", " ", 10)                 // Nro. Secuencial de Pago
	tmpLine += padLeft(" ", " ", 10)                 // cajero
	tmpLine += padLeft(CVSRow[2], " ", 8)            // Hora de pago Formato  HH:MM:SS
	tmpLine += padLeft(" ", " ", 2)                  // RELLENO
	tmpLine += padLeft("S", " ", 1)                  // Pago parcial S  -  Pago parcial             N  -  Pago total
	tmpLine += padLeft(importe, "0", 15)             // Impore pagado
	tmpLine += padLeft("20990908", " ", 8)           // Fecha de Vencimiento  (1ra., 2da. o 3ra.) del importe pagado
	tmpLine += padLeft("", "0", 8)                   // Número del cheque recibido, acreditado o rechazado
	tmpLine += padLeft(" ", " ", 3)                  // Banco emisor del cheque
	tmpLine += padLeft(" ", " ", 4)                  // Sucursal emisora del cheque
	tmpLine += padLeft(" ", " ", 7)                  // Codigo Postal
	tmpLine += padLeft(" ", " ", 8)                  // Fecha Presentación del Cheque
	tmpLine += padLeft(" ", " ", 11)                 // Cuit del emisor del cheque
	tmpLine += padLeft(" ", " ", 11)                 // Cuenta del emisor del cheque
	tmpLine += padLeft(" ", " ", 1)                  // Estado del cheque
	tmpLine += padLeft(" ", " ", 10)                 // Motivo de rechazo (Código Cimpra)/Motivo de Negociación ( Ver solapa "Codigo-Mot Resc-Rechazo" )
	tmpLine += padLeft(" ", " ", 2)                  // Canal de Pago
	tmpLine += padLeft(" ", " ", 10)                 // Nro. de  transacción de Pago
	tmpLine += padLeft(" ", " ", 10)                 // Fecha de ingreso  del cheque
	tmpLine += padLeft(" ", " ", 10)                 // Fecha de imputación del cheque
	tmpLine += padLeft(" ", " ", 13)                 // Espacios en blanco

	tmpLine += "\r\n"
	return tmpLine
}

func saveFile(textToSave string, DATPath string) {
	fileName := "BM17032020.0042.dat"

	f, err := os.Create(DATPath + fileName)
	check(err)
	l, err := f.WriteString(textToSave)
	check(err)
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func padLeft(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func padRight(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
