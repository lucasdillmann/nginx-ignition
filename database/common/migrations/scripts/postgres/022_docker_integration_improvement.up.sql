update host_route
set integration_option_id = concat(integration_option_id, ':host')
where integration_id in (
    select id
    from integration
    where driver = 'DOCKER'
);
