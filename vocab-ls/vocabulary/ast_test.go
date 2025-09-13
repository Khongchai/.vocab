package vocabulary

import (
	"context"
	"testing"
)

func TestInitializeAst(t *testing.T) {
	if NewAst(context.Background(), "", "", nil) == nil {
		t.Fatalf("Could not initialize ast, received nil")
	}
}

func TestValidSectionOpening(t *testing.T) {
	// timeText := "20/05/2025"
	// time, err := time.Parse(internal.DateRef, "20/05/2025")
	// if err != nil {
	// 	t.Errorf("Invalid time format")
	// }

	// dates := []*Date{
	// 	{
	// 		Text: fmt.Sprint(timeText),
	// 		Time: time,
	// 		TextRange: lsproto.Range{
	// 			Start: lsproto.Position{
	// 				Line:      0,
	// 				Character: 0,
	// 			},
	// 			End: lsproto.Position{
	// 				Line:      0,
	// 				Character: 10,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Text: fmt.Sprintf(" %s ", timeText),
	// 		Time: time,
	// 		TextRange: lsproto.Range{
	// 			Start: lsproto.Position{
	// 				Line:      0,
	// 				Character: 1,
	// 			},
	// 			End: lsproto.Position{
	// 				Line:      0,
	// 				Character: 11,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Text: fmt.Sprintf("## %s", timeText),
	// 		Time: time,
	// 		TextRange: lsproto.Range{
	// 			Start: lsproto.Position{
	// 				Line:      0,
	// 				Character: 3,
	// 			},
	// 			End: lsproto.Position{
	// 				Line:      0,
	// 				Character: 13,
	// 			},
	// 		},
	// 	},
	// }

	// // for i, date := range dates {
	// // ast := NewAst(context.Background())

	// // // act
	// // ast.Make("somewhere://overtherainbow", date.Text, nil)

	// // documents := ast.documents["somewhere://overtherainbow"]
	// // if documents == nil {
	// // 	t.Fatalf("Unable to parse %d, received nil", i)
	// // }
	// // if len(documents) == 0 {
	// // 	t.Fatalf("Parsed documents length is 0, expected 1")
	// // }
	// // if len(documents[0].Sections) == 0 {
	// // 	t.Fatalf("Parsed sections length is 0, expected 1")
	// // }

	// // result := documents[0].Sections[0].Date
	// // if *result != *date {
	// // 	t.Fatalf("Expected %+v, got %+v", *date, *result)
	// // }
	// // }
}

func TestValidMultipleSectionOpenings(t *testing.T) {
}
