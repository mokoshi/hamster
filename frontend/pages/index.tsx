import { Box, SimpleGrid } from '@chakra-ui/react';
import OpenOrders from './OpenOrders';
import Balance from './Balance';
import History from './History';

export default function Home() {
  return (
    <SimpleGrid p={2} columns={2} spacing={2}>
      <Box>
        <History />
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
