package all

import (
	// execute all `init` in packages
	_ "github.com/nktknshn/gomusiclibrary/cmd/database"
	_ "github.com/nktknshn/gomusiclibrary/cmd/library"
	_ "github.com/nktknshn/gomusiclibrary/cmd/scan"
)
