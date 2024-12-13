do $$
begin
    if (select count(1) from host where default_server) > 0 then
        return;
    end if;

    insert into host (
        id,
        enabled,
        default_server,
        websocket_support,
        http2_support,
        redirect_http_to_https,
        use_global_bindings
    ) values (
        '4614dc92-3680-49ad-92fe-f6043fe90703',
        true,
        true,
        true,
        true,
        false,
        false
    );

    insert into host_binding (
        id,
        host_id,
        type,
        ip,
        port
    ) values (
        'bd63060c-f903-4c45-854b-3f54e4670b9c',
        '4614dc92-3680-49ad-92fe-f6043fe90703',
        'HTTP',
        '0.0.0.0',
        80
    );

    insert into host_route (
        id,
        host_id,
        priority,
        type,
        source_path,
        static_response_code
    ) values (
        '5b87b3fa-e271-49fc-a9f4-ee329015f3f1',
        '4614dc92-3680-49ad-92fe-f6043fe90703',
        0,
        'STATIC_RESPONSE',
        '/',
        403
    );
end;
$$ language plpgsql;
