export function formatNumber(num: number) {
  return num.toLocaleString("id-ID", {
    style: "currency",
    currency: "IDR",
  });
}
