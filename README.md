<h1 align="center"><code>MacInfo Client</code></h1>

<div align="center">
  <sub>Created by <a href="https://github.com/jgengo">Jordane Gengo (Titus)</a></sub>
</div>
<div align="center">
  <sub>From <a href="https://hive.fi">Hive Helsinki</a></sub>
</div>
<div align="center">
    <sub>Highly inspired by <a href="#">maxreport</a> (max), itself inspired by <a href="#">macreport</a> (clem)</sub>
</div>

---

# WORK IN PROGRESS




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

