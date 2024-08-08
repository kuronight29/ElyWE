for _, videoPath := range videoPaths {
    // Khởi động MPV với video
    cmd := exec.Command("mpv", "--fs", "--loop", "--mute=yes", "--panscan=1.0", "--hwdec=auto", "--profile=low-latency", "--framedrop=no", "--scale=bilinear", "--dscale=bilinear", "--video-sync=display-resample", "--video-output-levels=full", videoPath)
    cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_CONSOLE, HideWindow: true}
    err := cmd.Start()
    if err != nil {
        showMessageBox("Error", "Failed to start MPV: \n"+err.Error(), MB_ICONERROR)
        return
    }

    // Chờ MPV window xuất hiện
    var mpvWindow uintptr
    var timeout = 0
    for {
        if timeout > 50 {
            showMessageBox("Error", "Failed to start MPV: Timeout (5s)\nYour video might be invalid.", MB_ICONERROR)
            return
        }
        mpvWindow, _, _ = procFindWindow.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("mpv"))), 0)
        if mpvWindow != 0 {
            break
        }
        timeout++
        time.Sleep(100 * time.Millisecond)
    }

    // Thiết lập kích thước và vị trí cho MPV window
    procSetWindowPos.Call(mpvWindow, 0, 0, 0, screenWidth, screenHeight, SWP_NOZORDER)
    procSetParent.Call(mpvWindow, workerw)

    // Chờ video kết thúc trước khi chuyển sang video tiếp theo
    time.Sleep(10 * time.Second) // Thay đổi thời gian chờ nếu cần
}
