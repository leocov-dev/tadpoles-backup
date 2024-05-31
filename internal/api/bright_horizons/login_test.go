package bright_horizons

import (
	"net/url"
	"tadpoles-backup/internal/interfaces"
	"tadpoles-backup/internal/test_utils"
	"testing"
)

func Test_fetchRequestVerificationToken(t *testing.T) {
	type args struct {
		client interfaces.HttpClient
		rvtUrl *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "parse RVT correctly",
			args: args{
				client: test_utils.NewMockClient(`
<!DOCTYPE html>
<html lang="en">
	<body>
		<head title="mocked"></head>
		<input name="my input" value="not the right value">
		<input name="__RequestVerificationToken" value="1234-AbCd_5678">
	</body>
</html>
`),
				rvtUrl: test_utils.MockUrl("https://mock.example.com"),
			},
			want: "1234-AbCd_5678",
		},
		{
			name: "RVT not in response",
			args: args{
				client: test_utils.NewMockClient(`
<!DOCTYPE html>
<html lang="en">
	<body>
		<head title="mocked"></head>
		Does not contain RVT
	</body>
</html>
`),
				rvtUrl: test_utils.MockUrl("https://mock.example.com"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchRequestVerificationToken(tt.args.client, tt.args.rvtUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchRequestVerificationToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchRequestVerificationToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startSaml(t *testing.T) {
	type args struct {
		client  interfaces.HttpClient
		samlUrl *url.URL
	}
	tests := []struct {
		name       string
		args       args
		wantAction string
		wantSaml   string
		wantErr    bool
	}{
		{
			name: "parse SAML correctly",
			args: args{
				client: test_utils.NewMockClient(`
<!DOCTYPE html>
<html lang="en">
	<body>
		<head title="mocked"></head>
		<input name="my input" value="not the right value">
		<form action="/action/path" method="POST">
		  <label for="fname">First name:</label><br>
		  <input type="text" id="fname" name="fname"><br>
		  <input type="text" id="samlId" name="SAMLResponse" value="ABC-1234_98ZZ">
		</form>
	</body>
</html>
`),
				samlUrl: test_utils.MockUrl("https://mock.example.com"),
			},
			wantAction: "/action/path",
			wantSaml:   "ABC-1234_98ZZ",
		},
		{
			name: "SAML not in page",
			args: args{
				client: test_utils.NewMockClient(`
<!DOCTYPE html>
<html lang="en">
	<body>
		<head title="mocked"></head>
		<input name="my input" value="not the right value">
		<form>
		  <label for="fname">First name:</label><br>
		  <input type="text" id="fname" name="fname"><br>
		</form>
	</body>
</html>
`),
				samlUrl: test_utils.MockUrl("https://mock.example.com"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAction, gotSaml, err := startSaml(tt.args.client, tt.args.samlUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("startSaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAction != tt.wantAction {
				t.Errorf("startSaml() got = %v, want %v", gotAction, tt.wantAction)
			}
			if gotSaml != tt.wantSaml {
				t.Errorf("startSaml() got1 = %v, want %v", gotSaml, tt.wantSaml)
			}
		})
	}
}
