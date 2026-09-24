update integration set parameters = json_patch(parameters, '{"legacyApi": true}') where driver = 'TRUENAS';
