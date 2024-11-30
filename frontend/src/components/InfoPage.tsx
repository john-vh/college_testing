import React from "react";
import AccountInfo from "./AccountInfo.tsx";
import { PostingInfoPage } from "./PostingInfo.tsx";
import { BusinessInfo } from "./BusinessInfo.tsx";
import { ApplicationInfo } from "./ApplicationInfo.tsx";

export enum PageTag {
    Account = "account",
    Application = "application",
    Business = "business",
    Approvals = "approvals",
    Postings = "postings"
}

interface InfoPageProps {
    page: PageTag
}

export const InfoPage = ({ page }: InfoPageProps) => {
    switch (page) {
        case PageTag.Account:
            return (
                <AccountInfo />
            );
        case PageTag.Application:
            return (
                <ApplicationInfo />
            );
        case PageTag.Business:
            return (
                <BusinessInfo />
            );
        case PageTag.Approvals:
            return (
                <AccountInfo />
            );
        case PageTag.Postings:
            return (
                <PostingInfoPage />
            );
        default:
            return (
                <AccountInfo />
            );
    }
}