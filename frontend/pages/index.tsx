import { Box, SimpleGrid } from '@chakra-ui/react';
import History from './history';
import OpenOrders from './OpenOrders';
import Balance from './Balance';

export default function Home() {
  return (
    <SimpleGrid p={2} columns={2} spacing={2}>
      <Box>
        <OpenOrders />
      </Box>
      <Box>
        <Balance />
      </Box>
      <Box>
        <History />
      </Box>
    </SimpleGrid>
  );
}
