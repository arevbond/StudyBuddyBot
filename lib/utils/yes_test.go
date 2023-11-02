package utils

import "testing"

func Test_YesOrNoCommand(t *testing.T) {
	testCases := []struct {
		text string
		want int
	}{
		{"да", IsYesCommand},
		{"ддддда", IsYesCommand},
		{"дддддаааааааа", IsYesCommand},
		{"дддддаааааааа.", IsYesCommand},
		{"     дддддаааааааа .", IsYesCommand},
		{"    дааа  .", IsYesCommand},
		{"    дааааа.  .", IsYesCommand},

		{"     да всем здарова пидоры .", UnsupportedCommand},
		{"     да всем здарова пидоры .", UnsupportedCommand},
		{"ддддддддд", UnsupportedCommand},
		{"аааааа", UnsupportedCommand},
		{"ф да залупа. .", UnsupportedCommand},
		{"   даб даб даб ", UnsupportedCommand},
		{" ", UnsupportedCommand},
		{"фывыфвфвнннннннннннеееетттттттт", UnsupportedCommand},

		{"нет", IsNoCommand},
		{"нет...", IsNoCommand},
		{"неееет", IsNoCommand},
		{"неееетттттттт", IsNoCommand},
		{"нннннннннннеееетттттттт", IsNoCommand},
		{"  фыфыв в  ы   нннннннннннеееетттттттт.", IsNoCommand},
	}

	for _, tc := range testCases {
		res := CheckYesOrNo(tc.text)
		if res != tc.want {
			t.Errorf("In \"%s\" resutl: %d, want: %d", tc.text, res, tc.want)
		}
	}
}
