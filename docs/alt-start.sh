#!/usr/bin/env bash

DIR="$( cd "$( dirname "$0" )" && pwd )"

CWD=$(pwd)

cd "$DIR/.."

mkdir -p demo/etc

export PYTHONPATH=$(readlink -f $(pwd))


echo $PYTHONPATH

for port in 8008 8009; do
    echo "Starting server on port $port... "

    https_port=$((port + 400))
    mkdir -p demo/$port
    pushd demo/$port

    #rm $DIR/etc/$port.config
    python3 -m synapse.app.homeserver \
        --generate-config \
        -H "localhost:$https_port" \
        --config-path "$DIR/etc/$port.config" \
        --report-stats no

    if ! grep -F "Customisation made by demo/start.sh" -q  $DIR/etc/$port.config; then
        printf '\n\n# Customisation made by demo/start.sh\n' >> $DIR/etc/$port.config

        echo "public_baseurl: http://localhost:$port/" >> $DIR/etc/$port.config

        echo 'enable_registration: true' >> $DIR/etc/$port.config

        # Warning, this heredoc depends on the interaction of tabs and spaces. Please don't
        # accidentaly bork me with your fancy settings.
		listeners=$(cat <<-PORTLISTENERS
		# Configure server to listen on both $https_port and $port
		# This overides some of the default settings above
		listeners:
		  - port: $https_port
		    type: http
		    tls: true
		    resources:
		      - names: [client, federation]

		  - port: $port
		    tls: false
		    bind_addresses: ['::1', '127.0.0.1']
		    type: http
		    x_forwarded: true
		    resources:
		      - names: [client, federation]
		        compress: false
		PORTLISTENERS
		)
        echo "${listeners}" >> $DIR/etc/$port.config

        # Disable tls for the servers
        printf '\n\n# Disable tls on the servers.' >> $DIR/etc/$port.config
        echo '# DO NOT USE IN PRODUCTION' >> $DIR/etc/$port.config
        echo 'use_insecure_ssl_client_just_for_testing_do_not_use: true' >> $DIR/etc/$port.config
        echo 'federation_verify_certificates: false' >> $DIR/etc/$port.config

        # Set tls paths
        echo "tls_certificate_path: \"$DIR/etc/localhost:$https_port.tls.crt\"" >> $DIR/etc/$port.config
        echo "tls_private_key_path: \"$DIR/etc/localhost:$https_port.tls.key\"" >> $DIR/etc/$port.config

        # Generate tls keys
        openssl req -x509 -newkey rsa:4096 -keyout $DIR/etc/localhost\:$https_port.tls.key -out $DIR/etc/localhost\:$https_port.tls.crt -days 365 -nodes -subj "/O=matrix"

        # Ignore keys from the trusted keys server
        echo '# Ignore keys from the trusted keys server' >> $DIR/etc/$port.config
        echo 'trusted_key_servers:' >> $DIR/etc/$port.config
        echo '  - server_name: "matrix.org"' >> $DIR/etc/$port.config
        echo '    accept_keys_insecurely: true' >> $DIR/etc/$port.config

        # Reduce the blacklist
        blacklist=$(cat <<-BLACK
		# Set the blacklist so that it doesn't include 127.0.0.1, ::1
		federation_ip_range_blacklist:
		  - '10.0.0.0/8'
		  - '172.16.0.0/12'
		  - '192.168.0.0/16'
		  - '100.64.0.0/10'
		  - '169.254.0.0/16'
		  - 'fe80::/64'
		  - 'fc00::/7'
		BLACK
		)
        echo "${blacklist}" >> $DIR/etc/$port.config
        echo "registration_shared_secret: \"secret\"" >> $DIR/etc/$port.config
    fi

    # Check script parameters
    if [ $# -eq 1 ]; then
        if [ $1 = "--no-rate-limit" ]; then

            # Disable any rate limiting
            ratelimiting=$(cat <<-RC
			rc_message:
			  per_second: 1000
			  burst_count: 1000
			rc_registration:
			  per_second: 1000
			  burst_count: 1000
			rc_login:
			  address:
			    per_second: 1000
			    burst_count: 1000
			  account:
			    per_second: 1000
			    burst_count: 1000
			  failed_attempts:
			    per_second: 1000
			    burst_count: 1000
			rc_admin_redaction:
			  per_second: 1000
			  burst_count: 1000
			rc_joins:
			  local:
			    per_second: 1000
			    burst_count: 1000
			  remote:
			    per_second: 1000
			    burst_count: 1000
			rc_3pid_validation:
			  per_second: 1000
			  burst_count: 1000
			rc_invites:
			  per_room:
			    per_second: 1000
			    burst_count: 1000
			  per_user:
			    per_second: 1000
			    burst_count: 1000
			RC
			)
            echo "${ratelimiting}" >> $DIR/etc/$port.config
        fi
    fi

    if ! grep -F "full_twisted_stacktraces" -q  $DIR/etc/$port.config; then
        echo "full_twisted_stacktraces: true" >> $DIR/etc/$port.config
    fi
    if ! grep -F "report_stats" -q  $DIR/etc/$port.config ; then
        echo "report_stats: false" >> $DIR/etc/$port.config
    fi

    python3 -m synapse.app.homeserver \
        --config-path "$DIR/etc/$port.config" \
        -D \

    popd
done

cd "$CWD"
