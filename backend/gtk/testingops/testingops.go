package testingops

import (
	"github.com/jwijenbergh/purego"
)

var (
	gskRendererRenderTexture func(uintptr, uintptr, uintptr) uintptr
	gdkTextureSaveToPng      func(uintptr, string) bool
	gObjectUnref             func(uintptr)
	gSignalEmitByName        func(uintptr, string)
)

func init() {
	gtk, err := purego.Dlopen("libgtk-4.so.1", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	purego.RegisterLibFunc(&gskRendererRenderTexture, gtk, "gsk_renderer_render_texture")
	purego.RegisterLibFunc(&gdkTextureSaveToPng, gtk, "gdk_texture_save_to_png")
	purego.RegisterLibFunc(&gObjectUnref, gtk, "g_object_unref")

	gobject, err := purego.Dlopen("libgobject-2.0.so.0", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	purego.RegisterLibFunc(&gSignalEmitByName, gobject, "g_signal_emit_by_name")
}

// RenderAndSave renders a GskRenderNode to a PNG file using the given
// GskRenderer. Passes NULL viewport so the node's own bounds are used.
func RenderAndSave(rendererPtr uintptr, nodePtr uintptr, filename string) bool {
	if rendererPtr == 0 || nodePtr == 0 {
		return false
	}
	texture := gskRendererRenderTexture(rendererPtr, nodePtr, 0)
	if texture == 0 {
		return false
	}
	defer gObjectUnref(texture)
	return gdkTextureSaveToPng(texture, filename)
}

// EmitClicked emits the "clicked" signal on a GObject.
func EmitClicked(widgetPtr uintptr) {
	if widgetPtr == 0 {
		return
	}
	gSignalEmitByName(widgetPtr, "clicked")
}
