package spell

import "fmt"

func (prescription *EditorialPrescription) Dump() {
	for i := 0; i < len(prescription.Actions); i++ {
		action := ""
		switch prescription.Actions[i] {
		case Insert:
			action = "I"
		case Delete:
			action = "D"
		case Replace:
			action = "R"
		case Match:
			action = "M"
		case Transposition:
			action = "T"
		case MissDouble:
			action = "U"
		case Duplicate:
			action = "P"
		case Triplet:
			action = "J"
		}
		fmt.Print(action, " ")
	}
	fmt.Println()

	for i := 0; i < len(prescription.Froms); i++ {
		fmt.Print(string(prescription.Froms[i]), " ")
	}
	fmt.Println()

	for i := 0; i < len(prescription.Tos); i++ {
		fmt.Print(string(prescription.Tos[i]), " ")
	}
	fmt.Println()
}
