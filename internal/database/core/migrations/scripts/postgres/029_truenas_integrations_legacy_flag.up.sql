update integration set parameters = parameters::jsonb || '{"legacyApi": true}'::jsonb where driver = 'TRUENAS';
