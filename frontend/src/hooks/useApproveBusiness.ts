import { useState, useEffect } from 'react';
import { PostingApplicationInfo } from './useApplicationInfo';
import { PostingInfo } from './useAllPostings';
import { Toaster, Position, Intent } from "@blueprintjs/core";

const AppToaster = Toaster.create({
  position: Position.BOTTOM_RIGHT,
});

export function useApproveBusiness() {

  const approveBusiness = (business_id: string) => {
    async function fetchData() {
      try {
        const response = await fetch(`${process.env.REACT_APP_API_URL}/admin/businesses/${business_id}/approve`, {
          method: 'POST',
          mode: 'cors',
          credentials: 'include',
        });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        AppToaster.show({
          message: "Business approved successfully!",
          intent: Intent.SUCCESS,
        });
      } catch (error) {
        console.log(error);
        AppToaster.show({
          message: "Failed to approve business.",
          intent: Intent.DANGER,
        });
      }
    }
    fetchData();
  };

  return approveBusiness;
}
