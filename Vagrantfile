# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV["LC_ALL"] = "en_US.UTF-8"

Vagrant.configure(2) do |config|

    config.vm.box = "ubuntu/xenial64"

    config.vm.define "pxe-pilot" do |server|
        config.vm.provider 'virtualbox' do |vb|
            vb.customize ['modifyvm', :id, '--memory', '2048']
            vb.customize ['modifyvm', :id, '--chipset', 'ich9']
        end

        server.vm.hostname = 'pxe-pilot'

        server.vm.network "private_network", ip: "10.69.70.70", mac: "000000000a70"
        server.vm.network "private_network", ip: "10.69.71.70", mac: "000000000b70"

        server.vm.provision "shell", privileged: false, path: "test/scripts/install-tools.sh"
        server.vm.provision "shell", privileged: false, path: "test/scripts/install-go.sh"
        server.vm.provision "shell", privileged: false, path: "test/scripts/install-bats.sh"
        server.vm.provision "shell", privileged: false, path: "test/scripts/install-mocks.sh"
        server.vm.provision "shell", privileged: false, path: "test/scripts/build.sh"
        server.vm.provision "shell", privileged: false, path: "test/scripts/install-service.sh"

        #
        # Running tests
        #
        server.vm.provision "shell", privileged: false, inline: <<-SHELL
            set -e

            echo "##############################################################"
            echo "### Running integration tests                              ###"
            echo "##############################################################"

            /vagrant/test/run.bats
        SHELL
    end
end
