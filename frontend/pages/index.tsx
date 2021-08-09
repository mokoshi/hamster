import { Box, SimpleGrid } from '@chakra-ui/react';
import OpenOrders from './OpenOrders';
import Balance from './Balance';
import OrderBooksHistory from './OrderBooksHistory';

export default function Home() {
  return (
    <SimpleGrid p={2} columns={2} spacing={2}>
      <Box>
        <OrderBooksHistory />
      </Box>
      <Box>
        <OpenOrders />
      </Box>
      <Box>
        <Balance />
      </Box>
    </SimpleGrid>
  );
}
