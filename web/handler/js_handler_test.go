package handler

import (
	"github.com/gofiber/fiber/v2"
	"testing"
)

func TestIcodeHandler(t *testing.T) {
	type args struct {
		fc *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	app := fiber.New()
	app.Group(ICoderHandler.Path, ICoderHandler.Handler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := icodeHandler(tt.args.fc); (err != nil) != tt.wantErr {
				t.Errorf("icodeHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
