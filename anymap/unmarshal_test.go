package anymap

import (
	"reflect"
	"testing"

	"github.com/myyrakle/caster/utils"
)

func TestUnmarshal(t *testing.T) {
	type shopOverview struct {
		PriceOfferNotConfirmed int  `json:"priceOfferNotConfirmed" bson:"priceOfferNotConfirmed"`
		ItemNotConfirmed       int  `json:"itemNotConfirmed" bson:"itemNotConfirmed"`
		ShippingNotProcessed   int  `json:"shippingNotProcessed" bson:"shippingNotProcessed"`
		ReturnNotProcessed     *int `json:"returnNotProcessed" bson:"returnNotProcessed"`
	}

	type args struct {
		source map[string]any
	}
	tests := []struct {
		name                    string
		args                    args
		mockDestinationVariable func() any
		want                    any
		wantErr                 bool
	}{
		{
			name: "정상 동작",
			args: args{
				source: map[string]any{
					"priceOfferNotConfirmed": 5,
					"itemNotConfirmed":       10,
					"shippingNotProcessed":   20,
					"returnNotProcessed":     30,
				},
			},
			want: shopOverview{
				PriceOfferNotConfirmed: 5,
				ItemNotConfirmed:       10,
				ShippingNotProcessed:   20,
				ReturnNotProcessed:     utils.ToPointer(30),
			},
			mockDestinationVariable: func() any {
				return shopOverview{}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variable := tt.mockDestinationVariable()

			if err := Unmarshal(tt.args.source, &variable); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(variable, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", variable, tt.want)
			}
		})
	}
}
