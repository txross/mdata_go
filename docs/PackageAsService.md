# Commands run to package golang as a service

1. Move environment files to correct location /etc/default<br>
    `sudo cp packaging/systemd/etc/default/sawtooth-mdata-tp-go /etc/default/sawtooth-mdata-tp-go`<br>
    `sudo chmod 644 /etc/default/sawtooth-mdata-tp-go`<br>

2. Move service files to correct location /lib/systemd/system<br>
    `sudo cp packaging/systemd/lib/systemd/system/sawtooth-mdata-tp-go.service /lib/systemd/system/`<br>
    `sudo chmod 644 /lib/systemd/system/sawtooth-mdata-go.service`<br>

3. Move binaries to correct location /usr/bin<br>
  - The processor
    `sudo cp sawtooth-mdata-tp-go /usr/bin/sawtooth-mdata-tp-go`<br>
    `sudo chmod 755 /usr/bin/sawtooth-mdata-tp-go`<br>

  - And the client
    `sudo cp mdata /usr/bin/mdata`<br>
    `sudo chmod 755 /usr/bin/mdata`<br>

4. Enable service <br>
    `sudo systemctl enable sawtooth-mdata-tp-go.service`<br>

5. Reload daemon <b>
    `sudo systemctl-reload`