update certificate
set provider_id = 'ACME',
    parameters = json_set(parameters, '$.certificateAuthority', 'LETSENCRYPT_PRODUCTION')
where provider_id = 'LETS_ENCRYPT';
