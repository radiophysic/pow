package utils

import (
    "bufio"
    "bytes"
    "io"
)

func LineCounter(r io.Reader) (int, error) {
    buf := make([]byte, 32*1024)
    count := 0
    lineSep := []byte{'\n'}

    for {
        c, err := r.Read(buf)
        count += bytes.Count(buf[:c], lineSep)

        switch {
        case err == io.EOF:
            return count, nil

        case err != nil:
            return count, err
        }
    }
}

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
    sc := bufio.NewScanner(r)
    for sc.Scan() {
        lastLine++
        if lastLine == lineNum {
            return sc.Text(), lastLine, sc.Err()
        }
    }
    return line, lastLine, io.EOF
}
