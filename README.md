# init



# notes for myself

- UUID is unique identifier for the Mac, it won't change
- hardware_serial and mac_address are changed when the mainboard is replaced.
-
-
# instructions

- copy `config/config.yml` into your `/etc`
```bash
cp config/config.yml /etc/macinfo.yml
```

or wherever you want but don't forget to specify -cfg argument while running the binary
```bash
./macinfo -cfg /opt/path/you/want/macinfo.yml
```

