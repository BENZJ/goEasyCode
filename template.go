package goEasyCode

const DoTempl = `
@Data
public class {{.ClassName}}DO{
    {{range .Columns}}
    /**
    * {{.ColumnComment}}
    */
    private {{.Type}} {{.ColumnName}}
    {{end}}
}

`

const DaoTempl = `
@Repository
public interface {{.ClassName}}Mapper {

    /**
     * 批量单条记录
     */
    Integer insert({{.ClassName}}DO item);

    /**
     * 批量插入
     */
    void batchInsert(List<{{.ClassName}}DO> list);

    /**
     * 根据id查询
     */
    {{.ClassName}}DO queryById(Long id);

    /**
     * 根据id查询
     */
    Integer update({{.ClassName}}DO item);

}
`

const MapperTempl = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="">

    <resultMap id="BaseResultMap" type="">
    {{range .Columns}}
      <result property="{{.ColumnName}}" column="{{.RealColumnName}}" jdbcType="{{.RealType}}"/>
    {{end}}
    </resultMap>

    <sql id="BaseColumnList">
        {{range  $key,$value := .Columns}}{{if (eq $key 0)}}{{- print $value.RealColumnName}}{{else}},{{- print $value.RealColumnName}}{{- end}}{{- end}}
    </sql>

    <!-- 批量插入-->
    <insert id="batchInsert" parameterType="java.util.List">
        INSERT INTO {{.TableName}} (
        {{range  $key,$value := .Columns}}{{if (eq $key 0)}}{{- print $value.RealColumnName}}{{else}},{{- print $value.RealColumnName}}{{- end}}{{- end}}
        )
        VALUES
        <foreach collection="list" item="item" index="index" separator=",">
            (
             {{range  $key,$value := .Columns}}{{if (eq $key 0)}}#{item.{{- print  $value.RealColumnName }}}{{else}},#{item.{{- print $value.RealColumnName}}}{{- end}}{{- end}}
            )
        </foreach>
    </insert>

    <!-- 插入-->
    <insert id="insert" keyProperty="id" useGeneratedKeys="true">
        insert into {{.TableName}}(
        {{range  $key,$value := .Columns}}{{if (eq $key 0)}}{{- print $value.RealColumnName}}{{else}},{{- print $value.RealColumnName}}{{- end}}{{- end}}
        )
        values (
             {{range  $key,$value := .Columns}}{{if (eq $key 0)}}#{item.{{- print  $value.RealColumnName }}}{{else}},#{item.{{- print $value.RealColumnName}}}{{- end}}{{- end}}
        )
    </insert>


    <!--通过主键更新-->
    <update id="update">
        update {{.TableName}}
        <set>
        {{- range  $key,$value := .Columns}}
        {{- if eq .Type "String" }}
            <if test="{{.RealColumnName}} != null and {{.RealColumnName}}  != '' ">
                {{ .RealColumnName}} = #{ {{- .RealColumnName}} },
            </if>
        {{- else}}
            <if test="{{.RealColumnName}} != null ">
                {{ .RealColumnName}} = #{ {{- .RealColumnName}} },
            </if>
        {{- end}}
        {{- end}}    
        </set>
        where ID = #{id}
    </update>

    <!--根据id查询-->
    <select id="queryById"  resultType="com.ctsec.org.wealth.dataobj.active.ActiveReportInfoDO">
        select <include refid="BaseColumnList"/>
        from {{.TableName}} where id = #{id}
    </select>

</mapper>
`
