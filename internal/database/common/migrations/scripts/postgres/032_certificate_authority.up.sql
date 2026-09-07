update certificate
set provider_id = 'ACME',
    parameters = jsonb_set(parameters::jsonb, '{certificateAuthority}', '"LETSENCRYPT_PRODUCTION"', true)::text
where provider_id = 'LETS_ENCRYPT';
