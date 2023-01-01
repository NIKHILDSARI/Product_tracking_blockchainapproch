package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Product_tracing struct {
	contractapi.Contract
}

type CollectedProduct_cert struct {
	Collect_location   string `json:"collect_location"`
	Density_reading    int    `json:"density_reading"`
	Lactometer_reading int    `json:"lactometer_reading"`
	Suppler_name       string `json:"suppler_name"`
	Collector_name     string `json:"collector_name"`
}

type Batch_cert struct {
	Type_of_milk           string   `json:"type_of_milk"`
	Manufacturing_loaction string   `json:"manufacturing_loaction"`
	Container_ids          []string `json:"container_ids"`
	Date_of_manufacture    string   `json:"date_of_manufacture"`
	Date_of_expiration     string   `json:"date_of_expiration"`
	Buyer                  string   `json:"buyer"`
	Cert_holder            string   `json:"cert_holder"`
}

type Shipping_cert struct {
	Vehicle_number    string   `json:"vehicle_number"`
	Gps_location      string   `json:"gps_location"`
	Batch_numbers     []string `json:"batch_numbers"`
	Shipping_agent_id string   `json:"shipping_agent_id"`
	Retailer_name     string   `json:"retailer_name"`
}

type Govt_Approvel_cert struct {
	Clearance_cert_id string `json:""`
	Sale_permit_id    string `json:""`
	Batch_number      string `json:""`
}

// *****************************************************************************************************************************************
func (p *Product_tracing) Generate_CollectedProduct_cert(
	ctx contractapi.TransactionContextInterface,
	location string,
	density_reading int,
	lacto_reading int,
	suppler string,
	collector string) (string, error) {

	collectedProduct_cert := CollectedProduct_cert{Collect_location: location,
		Density_reading:    density_reading,
		Lactometer_reading: lacto_reading,
		Suppler_name:       suppler,
		Collector_name:     collector,
	}

	collectedProduct_cert_btyes, err1 := json.Marshal(collectedProduct_cert)
	if err1 != nil {
		return "", err1
	}
	err1 = ctx.GetStub().PutState(get_containerID(), collectedProduct_cert_btyes)
	if err1 != nil {
		return "", err1
	}
	return get_containerID(), nil
}

func (p *Product_tracing) Generate_Batch_cert(
	ctx contractapi.TransactionContext,
	milk_type string,
	loaction string,
	container_ids []string,
	date_of_manf string,
	date_of_expi string,
	buyer string,
	cert_holder string,
	batch_id string,
	unit_id string) (Batch_cert, string, error) {

	batch_cert := Batch_cert{Type_of_milk: milk_type,
		Manufacturing_loaction: loaction,
		Container_ids:          container_ids,
		Date_of_manufacture:    date_of_manf,
		Date_of_expiration:     date_of_expi,
		Buyer:                  buyer,
		Cert_holder:            cert_holder,
	}

	B := Batch_cert{}
	batch_cert_bytes, err2 := json.Marshal(batch_cert)
	if err2 != nil {
		return B, "", err2
	}
	err2 = ctx.GetStub().PutState(batch_id, batch_cert_bytes)
	if err2 != nil {
		return B, "", err2
	}

	unit_id_Iotscanner := unit_id

	update_QuaryForConsumer(batch_cert, ctx, unit_id_Iotscanner, batch_id)

	return batch_cert, unit_id_Iotscanner, nil
}

func (p *Product_tracing) Generate_Shipping_cert(
	ctx contractapi.TransactionContext,
	vehicle_no string,
	gps_location string,
	batch_no []string,
	shipping_agent_id string,
	retailer_name string,
	unit_id string) (string, error) {

	shipping_cert := Shipping_cert{Vehicle_number: vehicle_no,
		Gps_location:      gps_location,
		Batch_numbers:     batch_no,
		Shipping_agent_id: shipping_agent_id,
		Retailer_name:     retailer_name,
	}
	cert_id := get_Shipping_ID(gps_location)
	shipping_cert_bytes, err3 := json.Marshal(shipping_cert)
	if err3 != nil {
		return "", err3
	}
	err3 = ctx.GetStub().PutState(cert_id, shipping_cert_bytes)
	if err3 != nil {
		return "", err3
	}
	update_QuaryForConsumer(shipping_cert, ctx, unit_id, cert_id)
	return cert_id, nil
}

func (p *Product_tracing) Gov_Approvel_cert(
	ctx contractapi.TransactionContext,
	clearance_id string,
	sale_permit_id string) (string, string, error) {

	Batch_serialnumber := get_Batch_serialNumber()

	unitid_for__Iotscanner := get_Queary_serialNumber()

	govt_Approvel_cert := Govt_Approvel_cert{Clearance_cert_id: clearance_id,
		Sale_permit_id: sale_permit_id,
		Batch_number:   Batch_serialnumber,
	}

	govt_Approvel_cert_bytes, err4 := json.Marshal(govt_Approvel_cert)
	if err4 != nil {
		return "", "", err4
	}
	err4 = ctx.GetStub().PutState(get_Approvel_ID(), govt_Approvel_cert_bytes)
	if err4 != nil {
		return "", "", err4
	}
	update_QuaryForConsumer(govt_Approvel_cert, ctx, unitid_for__Iotscanner, Batch_serialnumber)
	return Batch_serialnumber, unitid_for__Iotscanner, nil
}

// *******************************************************************************************************************************************
var Batch_Serialnumber int = 100
var Queary_Serialnumber int = 100

func get_Batch_serialNumber() string {
	Batch_Serialnumber += 1
	BatchSerialNumer_string := "DAIRYINDIA." + "BATCHID:" + strconv.Itoa(Batch_Serialnumber)
	return BatchSerialNumer_string
}
func get_Queary_serialNumber() string {
	Queary_Serialnumber += 1
	Quearyseriallnumber_string := "DAIRYINDIA." + "QUEARYID:" + strconv.Itoa(Queary_Serialnumber)
	return Quearyseriallnumber_string
}

// *********************************************************************************************************************************************
var Container_id int = 100

func get_containerID() string {
	Container_id += 1
	Container_id_string := "DAIRYINDIA." + "CONTAINERID:" + strconv.Itoa(Container_id)
	return Container_id_string
}

// ************************************************************************************************************************************************
var Shipment_ID int = 100
var Shipment_count int = 0

func get_Shipping_ID(location string) string {
	Shipment_ID += 1
	Shipment_count += 1
	Shipping_ID_string := "DAIRYINDIA." + "SHIPMENTID:" + strconv.Itoa(Shipment_ID) + ",STOP:" + strconv.Itoa(Shipment_count) + "," + location
	return Shipping_ID_string
}

// ********************************************************************************************************************************************************
var Approvel_ID int = 100

func get_Approvel_ID() string {
	Approvel_ID += 1
	Approvel_ID_string := "DAIRYINDIA." + "APPROVELID:" + strconv.Itoa(Approvel_ID)
	return Approvel_ID_string
}

// **********************************************************************************************************************************************************
// ************************************************************************************************************************************************************
type QuaryForConsumer struct {
	Batch_cert_Quary_ID         string
	Shipping_cert_Quary_ID      string
	Govt_Approvel_cert_Quary_ID string
	container_Quary_ids         []string
}

var quaryForConsumer QuaryForConsumer

func update_container_ids(
	ctx contractapi.TransactionContext,
	unit_ID string,
	cert_ID string,
	container_ids []string) (string, error) {

	ID_bytes, err5 := ctx.GetStub().GetState(unit_ID)
	if err5 != nil {
		return "", err5
	} else if ID_bytes != nil {
		return "already cert exist", nil

	} else {
		QuaryForConsumer_obj := QuaryForConsumer{}
		err5 = json.Unmarshal(ID_bytes, &QuaryForConsumer_obj)
		if err5 != nil {
			return "", err5
		}
		if QuaryForConsumer_obj.container_Quary_ids[0] != "" {
			return "already exist", nil
		}

		quaryForConsumer = QuaryForConsumer{container_Quary_ids: container_ids}
		quaryFromConsumer_bytes, err5 := json.Marshal(quaryForConsumer)
		if err5 != nil {
			return "", err5
		}
		err5 = ctx.GetStub().PutState(unit_ID, quaryFromConsumer_bytes)
		if err5 != nil {
			return "", err5
		}
		return "", nil
	}
}

func (s Batch_cert) update_QuaryForConsumer(
	ctx contractapi.TransactionContext,
	unit_ID string,
	cert_ID string) (string, error) {

	ID_bytes, err6 := ctx.GetStub().GetState(unit_ID)
	if err6 != nil {
		return "", err6
	} else if ID_bytes != nil {
		return "already cert exist", nil

	} else {

		QuaryForConsumer_obj := QuaryForConsumer{}
		err6 = json.Unmarshal(ID_bytes, &QuaryForConsumer_obj)
		if err6 != nil {
			return "", err6
		}

		if QuaryForConsumer_obj.Batch_cert_Quary_ID != "" {
			return "already exist", nil
		}

		quaryForConsumer = QuaryForConsumer{Batch_cert_Quary_ID: cert_ID}
		quaryFromConsumer_bytes, err6 := json.Marshal(quaryForConsumer)
		if err6 != nil {
			return "", err6
		}
		err6 = ctx.GetStub().PutState(unit_ID, quaryFromConsumer_bytes)
		if err6 != nil {
			return "", err6
		}
		update_container_ids(ctx, unit_ID, cert_ID, s.Container_ids)
		return "", nil
	}

}

func (s Shipping_cert) update_QuaryForConsumer(
	ctx contractapi.TransactionContext,
	unit_ID string,
	cert_ID string) (string, error) {

	ID_bytes, err7 := ctx.GetStub().GetState(unit_ID)
	if err7 != nil {
		return "", err7
	} else if ID_bytes != nil {
		return "already cert exist", nil

	} else {
		QuaryForConsumer_obj := QuaryForConsumer{}
		err7 = json.Unmarshal(ID_bytes, &QuaryForConsumer_obj)
		if err7 != nil {
			return "", err7
		}
		if QuaryForConsumer_obj.Shipping_cert_Quary_ID != "" {
			return "already exist", nil
		}

		quaryForConsumer = QuaryForConsumer{Shipping_cert_Quary_ID: cert_ID}
		quaryFromConsumer_bytes, err7 := json.Marshal(quaryForConsumer)
		if err7 != nil {
			return "", err7
		}
		err7 = ctx.GetStub().PutState(unit_ID, quaryFromConsumer_bytes)
		if err7 != nil {
			return "", err7
		}
		return "", nil
	}
}

func (s Govt_Approvel_cert) update_QuaryForConsumer(
	ctx contractapi.TransactionContext,
	unit_ID string,
	cert_ID string) (string, error) {

	ID_bytes, err8 := ctx.GetStub().GetState(unit_ID)
	if err8 != nil {
		return "", err8
	} else if ID_bytes != nil {
		return "already cert exist", nil

	} else {
		QuaryForConsumer_obj := QuaryForConsumer{}
		err8 = json.Unmarshal(ID_bytes, &QuaryForConsumer_obj)
		if err8 != nil {
			return "", err8
		}
		if QuaryForConsumer_obj.Govt_Approvel_cert_Quary_ID != "" {
			return "already exist", nil
		}

		quaryForConsumer = QuaryForConsumer{Govt_Approvel_cert_Quary_ID: cert_ID}
		quaryFromConsumer_bytes, err8 := json.Marshal(quaryForConsumer)
		if err8 != nil {
			return "", err8
		}
		err8 = ctx.GetStub().PutState(unit_ID, quaryFromConsumer_bytes)
		if err8 != nil {
			return "", err8
		}
		return "", nil
	}
}

type update_QuaryForConsumer_interface interface {
	update_QuaryForConsumer(

		ctx contractapi.TransactionContext,
		unit_ID string,
		cert_ID string) (string, error)
}

func update_QuaryForConsumer(
	cert update_QuaryForConsumer_interface,
	ctx contractapi.TransactionContext,
	unit_ID string,
	cert_ID string) {

	cert.update_QuaryForConsumer(ctx, unit_ID, cert_ID)
}

// *************************************************************************************************************************************************
func main() {
	cc := new(Product_tracing)
	chaincode, err := contractapi.NewChaincode(cc)
	if err != nil {
		fmt.Printf("err during newchaincode %v", err)
	}
	err = chaincode.Start()
	if err != nil {
		fmt.Printf("err while starting cc %v", err)
	}

}
