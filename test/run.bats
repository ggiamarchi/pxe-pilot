#!/usr/bin/env bats

load test_helper

@test "pxe-pilot host list with ARP cache not populated" {

    run pxe-pilot host list
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+-----------------------------------------------------------+-------------------+---------+-------------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION |                            MAC                            |     MGMT MAC      | MGMT IP | POWER STATE |" ]
    [ "${lines[2]}" = "+------+---------------+-----------------------------------------------------------+-------------------+---------+-------------+" ]
    [ "${lines[3]}" = "| h1   |               | 00:00:00:00:0b:01                                         | 00:00:00:00:0a:01 |         | Unknown     |" ]
    [ "${lines[4]}" = "| h2   |               | 00:00:00:00:0b:02                                         | 00:00:00:00:0a:02 |         | Unknown     |" ]
    [ "${lines[5]}" = "| h3   |               | 00:00:00:00:0b:03 | 00:00:00:00:0c:03 | 00:00:00:00:0d:03 |                   |         |             |" ]
    [ "${lines[6]}" = "| h4   |               | 00:00:00:00:0b:04                                         | 00:00:00:00:0a:04 |         | Unknown     |" ]
    [ "${lines[7]}" = "+------+---------------+-----------------------------------------------------------+-------------------+---------+-------------+" ]
}

@test "pxe-pilot host refresh & pxe-pilot host list" {

    run pxe-pilot host refresh
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+---------+" ]
    [ "${lines[1]}" = "| REFRESH |" ]
    [ "${lines[2]}" = "+---------+" ]
    [ "${lines[3]}" = "| OK      |" ]
    [ "${lines[4]}" = "+---------+" ]

    run pxe-pilot host list
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION |                            MAC                            |     MGMT MAC      |   MGMT IP   | POWER STATE |" ]
    [ "${lines[2]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
    [ "${lines[3]}" = "| h1   |               | 00:00:00:00:0b:01                                         | 00:00:00:00:0a:01 | 10.69.70.11 | Off         |" ]
    [ "${lines[4]}" = "| h2   |               | 00:00:00:00:0b:02                                         | 00:00:00:00:0a:02 | 10.69.70.12 | Off         |" ]
    [ "${lines[5]}" = "| h3   |               | 00:00:00:00:0b:03 | 00:00:00:00:0c:03 | 00:00:00:00:0d:03 |                   |             |             |" ]
    [ "${lines[6]}" = "| h4   |               | 00:00:00:00:0b:04                                         | 00:00:00:00:0a:04 | 10.69.70.14 | Off         |" ]
    [ "${lines[7]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
}

@test "pxe-pilot host on & off" {

    pxe-pilot host refresh

    run pxe-pilot host on h4
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+----+" ]
    [ "${lines[1]}" = "| NAME | ON |" ]
    [ "${lines[2]}" = "+------+----+" ]
    [ "${lines[3]}" = "| h4   | OK |" ]
    [ "${lines[4]}" = "+------+----+" ]

    run pxe-pilot host list
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[6]}" = "| h4   |               | 00:00:00:00:0b:04                                         | 00:00:00:00:0a:04 | 10.69.70.14 | On          |" ]

    run pxe-pilot host off h4
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+-----+" ]
    [ "${lines[1]}" = "| NAME | OFF |" ]
    [ "${lines[2]}" = "+------+-----+" ]
    [ "${lines[3]}" = "| h4   | OK  |" ]
    [ "${lines[4]}" = "+------+-----+" ]

    run pxe-pilot host list
    [ "${status}" -eq 0 ]
    [ "${lines[6]}" = "| h4   |               | 00:00:00:00:0b:04                                         | 00:00:00:00:0a:04 | 10.69.70.14 | Off         |" ]
}

@test "pxe-pilot bootloader list" {

    run pxe-pilot bootloader list
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+----------+-----------------------+---------------+" ]
    [ "${lines[1]}" = "|   NAME   |         FILE          |  CONFIG PATH  |" ]
    [ "${lines[2]}" = "+----------+-----------------------+---------------+" ]
    [ "${lines[3]}" = "| pxelinux | pxelinux.0            |               |" ]
    [ "${lines[4]}" = "| syslinux | syslinux.efi          |               |" ]
    [ "${lines[5]}" = "| grub     | grubnetx64.efi.signed | grub/grub.cfg |" ]
    [ "${lines[6]}" = "+----------+-----------------------+---------------+" ]
}

@test "pxe-pilot config list" {

    run pxe-pilot config list
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" =  "+-------+------------+" ]
    [ "${lines[1]}" =  "| NAME  | BOOTLOADER |" ]
    [ "${lines[2]}" =  "+-------+------------+" ]
    [ "${lines[3]}" =  "| local | pxelinux   |" ]
    [ "${lines[4]}" =  "| sheep | syslinux   |" ]
    [ "${lines[5]}" =  "| a     | grub       |" ]
    [ "${lines[6]}" =  "| b     | grub       |" ]
    [ "${lines[7]}" =  "| c     | grub       |" ]
    [ "${lines[8]}" =  "| grml  | grub       |" ]
    [ "${lines[9]}" =  "| z     | grub       |" ]
    [ "${lines[10]}" = "+-------+------------+" ]
}

@test "pxe-pilot config show" {

    run pxe-pilot config show sheep
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "default sheeplive" ]
    [ "${lines[1]}" = "label sheeplive" ]
    [ "${lines[2]}" = "    kernel /sheep-live/vmlinuz" ]
    [ "${lines[3]}" = "    append boot=live fetch=http://1.2.3.4/sheep-live.squashfs initrd=/sheep-live/initrd.img ssh=sheep console=ttyS1,57600n8 sheep.script=http://1.2.3.4/sheep sheep.config=http://1.2.3.4/config.yml" ]
}

@test "pxe-pilot config deploy" {

    pxe-pilot host refresh

    run pxe-pilot config deploy sheep h1
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+----------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION | REBOOTED |" ]
    [ "${lines[2]}" = "+------+---------------+----------+" ]
    [ "${lines[3]}" = "| h1   | sheep         | NO       |" ]
    [ "${lines[4]}" = "+------+---------------+----------+" ]

    run pxe-pilot host list
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION |                            MAC                            |     MGMT MAC      |   MGMT IP   | POWER STATE |" ]
    [ "${lines[2]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
    [ "${lines[3]}" = "| h1   | sheep         | 00:00:00:00:0b:01                                         | 00:00:00:00:0a:01 | 10.69.70.11 | Off         |" ]
    [ "${lines[4]}" = "| h2   |               | 00:00:00:00:0b:02                                         | 00:00:00:00:0a:02 | 10.69.70.12 | Off         |" ]
    [ "${lines[5]}" = "| h3   |               | 00:00:00:00:0b:03 | 00:00:00:00:0c:03 | 00:00:00:00:0d:03 |                   |             |             |" ]
    [ "${lines[6]}" = "| h4   |               | 00:00:00:00:0b:04                                         | 00:00:00:00:0a:04 | 10.69.70.14 | Off         |" ]
    [ "${lines[7]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]

    [ -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-01 ]
    [ "$(readlink ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-01)" = "${TFTP_ROOT}/pxelinux.cfg/conf/syslinux/sheep" ]

    rm -f ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-01
    [ ! -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-01 ]

    run pxe-pilot host list
    [ "${status}" -eq 0 ]
    [ "${lines[3]}" = "| h1   |               | 00:00:00:00:0b:01                                         | 00:00:00:00:0a:01 | 10.69.70.11 | Off         |" ]
}

@test "pxe-pilot config deploy --now" {

    pxe-pilot host refresh

    run pxe-pilot config deploy --now sheep h1
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+----------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION | REBOOTED |" ]
    [ "${lines[2]}" = "+------+---------------+----------+" ]
    [ "${lines[3]}" = "| h1   | sheep         | YES      |" ]
    [ "${lines[4]}" = "+------+---------------+----------+" ]

    run pxe-pilot host list
    [ "${status}" -eq 0 ]
    [ "${lines[3]}" = "| h1   | sheep         | 00:00:00:00:0b:01                                         | 00:00:00:00:0a:01 | 10.69.70.11 | On          |" ]
}

@test "pxe-pilot config deploy / with multiple mac addresses" {

    pxe-pilot host refresh

    run pxe-pilot config deploy local h3
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+----------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION | REBOOTED |" ]
    [ "${lines[2]}" = "+------+---------------+----------+" ]
    [ "${lines[3]}" = "| h3   | local         | NO       |" ]
    [ "${lines[4]}" = "+------+---------------+----------+" ]

    run pxe-pilot host list
    flush_log
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
    [ "${lines[1]}" = "| NAME | CONFIGURATION |                            MAC                            |     MGMT MAC      |   MGMT IP   | POWER STATE |" ]
    [ "${lines[2]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]
    [ "${lines[3]}" = "| h1   |               | 00:00:00:00:0b:01                                         | 00:00:00:00:0a:01 | 10.69.70.11 | Off         |" ]
    [ "${lines[4]}" = "| h2   |               | 00:00:00:00:0b:02                                         | 00:00:00:00:0a:02 | 10.69.70.12 | Off         |" ]
    [ "${lines[5]}" = "| h3   | local         | 00:00:00:00:0b:03 | 00:00:00:00:0c:03 | 00:00:00:00:0d:03 |                   |             |             |" ]
    [ "${lines[6]}" = "| h4   |               | 00:00:00:00:0b:04                                         | 00:00:00:00:0a:04 | 10.69.70.14 | Off         |" ]
    [ "${lines[7]}" = "+------+---------------+-----------------------------------------------------------+-------------------+-------------+-------------+" ]

    [ -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-03 ]
    [ "$(readlink ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-03)" = "${TFTP_ROOT}/pxelinux.cfg/conf/pxelinux/local" ]

    [ -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0c-03 ]
    [ "$(readlink ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0c-03)" = "${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-03" ]

    [ -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0d-03 ]
    [ "$(readlink ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0d-03)" = "${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-03" ]

    rm -f ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-03 ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0c-03 ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0d-03
    [ ! -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0b-03 ]
    [ ! -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0c-03 ]
    [ ! -e ${TFTP_ROOT}/pxelinux.cfg/01-00-00-00-00-0d-03 ]

    run pxe-pilot host list
    [ "${status}" -eq 0 ]
    [ "${lines[5]}" = "| h3   |               | 00:00:00:00:0b:03 | 00:00:00:00:0c:03 | 00:00:00:00:0d:03 |                   |             |             |" ]
}
