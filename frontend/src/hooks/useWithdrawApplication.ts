interface WithdrawApplicationProps {
  business_id: string,
  post_id: number,
  user_id: string
}

export function useWithdrawApplication() {

  const withdrawApplication = ({ business_id, post_id, user_id }: WithdrawApplicationProps) => {
    async function fetchData() {
      try {
        const response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/posts/${post_id}/applications/${user_id}/withdraw`, { method: "POST", mode: "cors", credentials: 'include' });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
      } catch (error) {
        console.log(error);
      }
    }
    fetchData();
  }

  return withdrawApplication;
}
