import memgraph from 'neo4j-driver';
import type { Driver } from 'neo4j-driver';
import { MEMEGRAPH_URI, MEMGRAPH_USER, MEMGRAPH_PASSWORD } from '$env/static/private';

export const driverInstance: Driver = memgraph.driver(
    MEMEGRAPH_URI || 'bolt://localhost:7687',
    memgraph.auth.basic(
        MEMGRAPH_USER || 'memgraph',
        MEMGRAPH_PASSWORD || 'memgraph'
    )
);
