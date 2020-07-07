package pacemaker

type Pacemaker interface {
	UpdateHighQC()
	OnBeat()
	OnNextSyncView()
	OnReceiverNewView()
}

type pacemakerImpl struct {

}
