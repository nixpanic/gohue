/*
* scene_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

package hue

import (
    "testing"
    "fmt"
)

func TestGetScenes(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128", "427de8bd6d49f149c8398e4fc08f")
    scenes, _ := bridge.GetScenes()
    for scene := range scenes {
        fmt.Println("SCENE: ", scenes[scene])
    }

    individual, _ := bridge.GetScene(scenes[0].ID)
    fmt.Println("Individual scene: ", individual)
}