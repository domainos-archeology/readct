# readct

A CLI for showing the index of apollo rbak format tape archives (such as those found on [bitsavers](http://bitsavers.org/bits/Apollo/SR10.4/), and for extracting files from those archives.

I got kinda tired of dealing with the `rbak` source from the apollo archives, and also getting it to compile on macos.  So, I rewrote it in go.

examples:

### index from a collection of tape archive files

```
% ./readct \
  017286-001.CRTG_STD_SFW_BOOT_1-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct \
  017287-002.CRTG_STD_SFW_2-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct \
  017287-001.CRTG_STD_SFW_1-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct \
  017287-003.CRTG_STD_SFW_3-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct
read_ct: reading from ../../sr10.2/017286-001.CRTG_STD_SFW_BOOT_1-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct
(dir) bscom
(file) bscom/rbak_shell  (nil 306009)
(dir) sau7
(file) sau7/calendar  (nil 154138)
(file) sau7/chuvol  (nil 279066)
(file) sau7/color7_ee_write  (nil 164378)
(file) sau7/config  (nil 154138)
(file) sau7/dex  (nil 212948)
(file) sau7/disp7e.dex  (nil 80724)
(file) sau7/domain_os  (nil 766696)
(file) sau7/domain_os.map  (unstruct 99226)
(file) sau7/fbs  (nil 152090)
(file) sau7/invol  (nil 186906)
(file) sau7/rwvol  (nil 147994)
(file) sau7/salvol  (nil 313882)
(file) sau7/self_test  (nil 13176) 
...
etc
```

### extracting from a collection of tape archive files

```
% ./readct x \
  017286-001.CRTG_STD_SFW_BOOT_1-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct \
  017287-002.CRTG_STD_SFW_2-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct \
  017287-001.CRTG_STD_SFW_1-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct \
  017287-003.CRTG_STD_SFW_3-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct
read_ct: reading from ../../sr10.2/017286-001.CRTG_STD_SFW_BOOT_1-REV.00-SYSTEM_RELEASE_VER_SR10.2-RAI.ct
(dir) bscom
(file) bscom/rbak_shell  (nil 306009)
(dir) sau7
(file) sau7/calendar  (nil 154138)
(file) sau7/chuvol  (nil 279066)
(file) sau7/color7_ee_write  (nil 164378)
(file) sau7/config  (nil 154138)
(file) sau7/dex  (nil 212948)
(file) sau7/disp7e.dex  (nil 80724)
(file) sau7/domain_os  (nil 766696)
(file) sau7/domain_os.map  (unstruct 99226)
(file) sau7/fbs  (nil 152090)
(file) sau7/invol  (nil 186906)
(file) sau7/rwvol  (nil 147994)
(file) sau7/salvol  (nil 313882)
(file) sau7/self_test  (nil 13176) 
...
etc

% ls
broken     com/       etc/       lib/       sau6/      sau8/      sys/       tmp/       usr/
bscom/     dev/       install/   sau5/      sau7/      sau9/      sysboot    user_data/
```
