import { Button, Card, FormGroup, H2, H3, InputGroup, Intent, Label, Tag } from "@blueprintjs/core";
import { Toaster, Position } from "@blueprintjs/core";
import React, { useMemo } from "react";
import { UserApplicationInfo, useUserApplicationInfo } from "../hooks/useUserApplicationInfo.ts";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import { useWithdrawApplication } from "../hooks/useWithdrawApplication.ts";
import AccountInfo from "./AccountInfo.tsx";

const toaster = Toaster.create({
    position: Position.TOP,
});

export const UserApplicationInfoPage: React.FC = () => {
    const { applicationInfo, fetchData } = useUserApplicationInfo();
    const account = useAccountInfo();
    const withdrawApplication = useWithdrawApplication();

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const sortedData = useMemo(
        () =>
            [...applicationInfo]
                .filter((entry) => entry.status !== "withdrawn")
                .sort((a, b) => -a.created_at.localeCompare(b.created_at)),
        [applicationInfo]
    );

    if (account == null) {
        return <AccountInfo />
    }

    const handleWithdraw = async (application: UserApplicationInfo) => {
        const { business: { id: business_id }, post: { id: post_id } } = application;
        const user_id = account?.id;
        withdrawApplication({ business_id, post_id, user_id });
        await delay(100);
        await fetchData();
    }

    return (
        <div className="content">
            <H2 style={{ marginBottom: "0px" }}>Manage Applications</H2>
            {sortedData.map((application) => (
                <Card interactive={true}>
                    <div className="Flex" style={{ justifyContent: "space-between" }}>
                        <H3>Application Info</H3>
                        <Tag round intent={
                            application.status === "active"
                                ? Intent.SUCCESS
                                : application.status === "pending"
                                    ? Intent.WARNING
                                    : Intent.DANGER
                        }>{application.status.toLocaleUpperCase()}</Tag>
                    </div>

                    <Tag minimal style={{ marginBottom: "7px" }}>{formatDate(application.created_at)}</Tag>
                    <FormGroup label="Title" labelFor="title">
                        <InputGroup id="title" defaultValue={application.post.title}></InputGroup>
                    </FormGroup>
                    <FormGroup label="Business"
                        labelFor="business" >
                        <InputGroup id="business" defaultValue={application.business.name} readOnly />
                    </FormGroup>
                    <FormGroup label="Compensation"
                        labelFor="pay" >
                        <InputGroup id="pay" defaultValue={"$" + application.post.pay.toString()} readOnly />
                    </FormGroup>
                    <div className="Flex" style={{ justifyContent: "space-between" }}>
                        {application.status !== "withdrawn" && <Button onClick={() => handleWithdraw(application)}>Withdraw application</Button>}
                        <Button rightIcon="issue" minimal />
                    </div>
                </Card>
            ))}
        </div>
    );
};


export function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
        dateStyle: 'medium', // Example: Nov 15, 2024
        timeStyle: 'short',  // Example: 9:28 PM
    }).format(date);
};