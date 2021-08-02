import { Box, SimpleGrid } from '@chakra-ui/react';
import { useOrderBookHistoriesQuery } from '../model/orderBooksHistory';
import { Suspense } from 'react';
import History from './history';

export default function Home() {
  const result = useOrderBookHistoriesQuery('2021-08-01', '2021-08-02');
  console.log(result.data);

  return (
    <SimpleGrid columns={2} spacing={10}>
      <Box>
        <Suspense fallback={<div>loading</div>}>
          <History />
        </Suspense>
      </Box>
    </SimpleGrid>
  );
}