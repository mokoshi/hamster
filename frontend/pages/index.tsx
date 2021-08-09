import { Box, SimpleGrid } from '@chakra-ui/react';
import History from './history';
import OpenOrders from './openOrders';

export default function Home() {
  return (
    <SimpleGrid columns={2} spacing={10}>
      <Box>
        <OpenOrders />
      </Box>
      <Box>
        <History />
      </Box>
    </SimpleGrid>
  );
}
