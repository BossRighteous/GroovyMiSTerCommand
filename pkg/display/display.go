package display

import (
	"fmt"
	"time"
)

var TickDuration time.Duration = time.Duration(int64(1000000000 / 60))
var TimerDuration time.Duration = time.Duration(5 * time.Second)

type MiSTerDisplay struct {
	IsRunning bool
	Client    *UdpDisplayClient
	StartChan chan bool
	StopChan  chan bool
	Ticker    *time.Ticker
	Timer     *time.Timer
	Frame     *BGR8
}

func (disp *MiSTerDisplay) SafeClose() {
	fmt.Println("Display: SafeClose")
	if disp.IsRunning {
		fmt.Println("Display: Running close seq")
		disp.Ticker.Stop()
		disp.Timer.Stop()
		disp.IsRunning = false
		disp.Client.CmdClose()
		time.Sleep(time.Millisecond * 250)
	}
}

func (disp *MiSTerDisplay) SafeOpen() {
	fmt.Println("Display: SafeOpen")
	if !disp.IsRunning {
		fmt.Println("Display: Running open seq")
		disp.IsRunning = true
		disp.Client.Open()
		disp.Client.CmdInit()
		disp.Client.CmdSwitchres()
		disp.StartChan <- true
		disp.Ticker.Reset(TickDuration)
	}
}

func (disp *MiSTerDisplay) BlitText(txt []string) {
	fmt.Println("Display: starting text broadcast for 5s")
	disp.SafeOpen()
	disp.Timer.Reset(TimerDuration)
	disp.Frame = TextToBGR8(txt)
}

func NewMiSTerDisplay(host string) *MiSTerDisplay {
	disp := &MiSTerDisplay{
		StopChan:  make(chan bool),
		StartChan: make(chan bool),
		IsRunning: false,
		Client:    NewUdpClient(host),
		Frame:     TextToBGR8([]string{}),
		Timer:     time.NewTimer(TimerDuration),
		Ticker:    time.NewTicker(TickDuration),
	}
	disp.Timer.Stop()
	disp.Ticker.Stop()

	go func() {
		for {
			select {
			case <-disp.StartChan:
			case <-disp.StopChan:
				disp.SafeClose()
			case <-disp.Timer.C:
				disp.SafeClose()
			case <-disp.Ticker.C:
				if disp.IsRunning {
					disp.Client.CmdBlit(disp.Frame.Pix)
				}
			}
		}
	}()

	return disp
}
