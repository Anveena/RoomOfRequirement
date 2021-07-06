package ezNetworking

type EZNetError struct {
	IsSettingError  bool
	IsCheckingError bool
	IsLoadingError  bool
	errorDesc       string
}

func (e *EZNetError) Error() string {
	return e.errorDesc
}
