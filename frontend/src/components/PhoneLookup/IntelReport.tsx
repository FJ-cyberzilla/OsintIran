// frontend/src/components/PhoneLookup/IntelReport.tsx
interface IntelReportProps {
  data: PhoneIntelData;
}

export const IntelReport: React.FC<IntelReportProps> = ({ data }) => {
  return (
    <div className="intel-report">
      <div className="report-header">
        <h2>گزارش هوشمند شماره: {data.phoneNumber}</h2>
        <div className="confidence-score">
          اطمینان: <span className={`score ${getScoreClass(data.confidence)}`}>
            {data.confidence}%
          </span>
        </div>
      </div>

      <div className="report-sections">
        {/* Basic Information */}
        <Section title="اطلاعات پایه">
          <InfoGrid>
            <InfoItem label="اپراتور" value={data.operator} />
            <InfoItem label="منطقه" value={data.region} />
            <InfoItem label="نوع خط" value={data.lineType} />
            <InfoItem label="وضعیت" value={data.status} />
          </InfoGrid>
        </Section>

        {/* Social Media Profiles */}
        <Section title="حساب‌های شبکه‌های اجتماعی">
          <SocialProfiles profiles={data.socialProfiles} />
        </Section>

        {/* Online Presence */}
        <Section title="فعالیت آنلاین">
          <OnlinePresence data={data.onlinePresence} />
        </Section>

        {/* Risk Assessment */}
        <Section title="ارزیابی امنیتی">
          <RiskAssessment risks={data.riskAssessment} />
        </Section>

        {/* Timeline */}
        <Section title="گاهشمار فعالیت">
          <ActivityTimeline activities={data.activities} />
        </Section>
      </div>

      <div className="report-actions">
        <button className="btn-primary">ذخیره گزارش</button>
        <button className="btn-secondary">چاپ گزارش</button>
        <button className="btn-warning">گزارش تخلف</button>
      </div>
    </div>
  );
};
