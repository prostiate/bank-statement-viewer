import { useQuery } from "@tanstack/react-query";
import { formatNumber } from "../utils/formatNumber";

type BalanceData = {
  balance: number;
};

type BalanceResponse = {
  data: BalanceData;
};

function GetBalance() {
  const { data, isPending, isError, error } = useQuery<BalanceResponse>({
    queryKey: ["balance"],
    queryFn: () => fetch("/api/balance").then((r) => r.json()),
  });

  if (isPending) return <div>Loading balance...</div>;
  if (isError) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>Balance</h2>
      <p>{formatNumber(data.data.balance)}</p>
    </div>
  );
}

export default GetBalance;
