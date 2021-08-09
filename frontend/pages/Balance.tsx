import { Table, Tr, Th, Td, Text, AlertIcon, Alert, Box, Skeleton, Tbody } from '@chakra-ui/react';
import { useBalanceQuery } from '../model/balance';
import React from 'react';
import Card from '../components/Card';

const Balance: React.FC = () => {
  return (
    <Card title='残高'>
      <BalanceContent />
    </Card>
  );
};

const BalanceContent: React.FC = () => {
  const { data: balance, isLoading, isError } = useBalanceQuery();

  if (isError) {
    return (
      <Alert status='error'>
        <AlertIcon />
        エラーが発生しました
      </Alert>
    );
  }
  if (isLoading || !balance) {
    return <Skeleton height={40} />;
  }

  return (
    <Table variant='simple'>
      <Tbody>
        <Tr>
          <Th>jpy</Th>
          <Td isNumeric>{balance.jpy}</Td>
        </Tr>
        <Tr>
          <Th>btc</Th>
          <Td isNumeric>{balance.btc}</Td>
        </Tr>
        <Tr>
          <Th>jpyReserved</Th>
          <Td isNumeric>{balance.jpyReserved}</Td>
        </Tr>
        <Tr>
          <Th>btcReserved</Th>
          <Td isNumeric>{balance.btcReserved}</Td>
        </Tr>
        <Tr>
          <Th>jpyLendInUse</Th>
          <Td isNumeric>{balance.jpyLendInUse}</Td>
        </Tr>
        <Tr>
          <Th>btcLendInUse</Th>
          <Td isNumeric>{balance.btcLendInUse}</Td>
        </Tr>
        <Tr>
          <Th>jpyLent</Th>
          <Td isNumeric>{balance.jpyLent}</Td>
        </Tr>
        <Tr>
          <Th>btcLent</Th>
          <Td isNumeric>{balance.btcLent}</Td>
        </Tr>
        <Tr>
          <Th>jpyDebt</Th>
          <Td isNumeric>{balance.jpyDebt}</Td>
        </Tr>
        <Tr>
          <Th>btcDebt</Th>
          <Td isNumeric>{balance.btcDebt}</Td>
        </Tr>
      </Tbody>
    </Table>
  );
};

export default Balance;
